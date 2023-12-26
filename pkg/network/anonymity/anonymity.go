package anonymity

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/anonymity/queue"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/state"

	"github.com/number571/go-peer/pkg/network/anonymity/adapters"
	anon_logger "github.com/number571/go-peer/pkg/network/anonymity/logger"
	net_message "github.com/number571/go-peer/pkg/network/message"
)

var (
	_ INode = &sNode{}
)

type sNode struct {
	fMutex         sync.Mutex
	fState         state.IState
	fSettings      ISettings
	fLogger        logger.ILogger
	fWrapperDB     IWrapperDB
	fNetwork       network.INode
	fQueue         queue.IMessageQueue
	fFriends       asymmetric.IListPubKeys
	fHandleRoutes  map[uint32]IHandlerF
	fHandleActions map[string]chan []byte
}

func NewNode(
	pSett ISettings,
	pLogger logger.ILogger,
	pWrapperDB IWrapperDB,
	pNetwork network.INode,
	pQueue queue.IMessageQueue,
	pFriends asymmetric.IListPubKeys,
) INode {
	return &sNode{
		fState:         state.NewBoolState(),
		fSettings:      pSett,
		fLogger:        pLogger,
		fWrapperDB:     pWrapperDB,
		fNetwork:       pNetwork,
		fQueue:         pQueue,
		fFriends:       pFriends,
		fHandleRoutes:  make(map[uint32]IHandlerF),
		fHandleActions: make(map[string]chan []byte),
	}
}

func (p *sNode) Run(pCtx context.Context) error {
	enableFunc := func() error {
		p.fNetwork.HandleFunc(
			p.fSettings.GetNetworkMask(),
			p.handleWrapper(),
		)
		return nil
	}
	if err := p.fState.Enable(enableFunc); err != nil {
		return fmt.Errorf("node running error: %w", err)
	}

	defer func() {
		disableFunc := func() error {
			p.fNetwork.HandleFunc(p.fSettings.GetNetworkMask(), nil)
			return nil
		}
		if err := p.fState.Disable(disableFunc); err != nil {
			panic(err)
		}
	}()

	chErr := make(chan error)
	go func() { chErr <- p.fQueue.Run(pCtx) }()

	for {
		select {
		case <-pCtx.Done():
			return <-chErr
		default:
			msg := p.fQueue.DequeueMessage(pCtx)
			if msg == nil {
				break
			}

			logBuilder := anon_logger.NewLogBuilder(p.fSettings.GetServiceName())
			if ok, _ := p.storeHashWithBroadcast(pCtx, logBuilder, msg); !ok {
				// internal logger
				break
			}

			// enrich logger
			logBuilder.
				WithPubKey(p.fQueue.GetClient().GetPubKey()).
				WithType(anon_logger.CLogBaseBroadcast)

			p.fLogger.PushInfo(logBuilder)
		}
	}
}

func (p *sNode) GetLogger() logger.ILogger {
	return p.fLogger
}

func (p *sNode) GetSettings() ISettings {
	return p.fSettings
}

func (p *sNode) GetWrapperDB() IWrapperDB {
	return p.fWrapperDB
}

func (p *sNode) GetNetworkNode() network.INode {
	return p.fNetwork
}

func (p *sNode) GetMessageQueue() queue.IMessageQueue {
	return p.fQueue
}

// Return f2f structure.
func (p *sNode) GetListPubKeys() asymmetric.IListPubKeys {
	return p.fFriends
}

func (p *sNode) HandleFunc(pHead uint32, pHandle IHandlerF) INode {
	p.setRoute(pHead, pHandle)
	return p
}

// Send message without response waiting.
func (p *sNode) BroadcastPayload(pCtx context.Context, pRecv asymmetric.IPubKey, pPld adapters.IPayload) error {
	// internal logger
	if err := p.enqueuePayload(pCtx, cIsRequest, pRecv, pPld.ToOrigin()); err != nil {
		return fmt.Errorf("broadcast payload: %w", err)
	}
	return nil
}

// Send message with response waiting.
// Payload head must be uint32.
func (p *sNode) FetchPayload(pCtx context.Context, pRecv asymmetric.IPubKey, pPld adapters.IPayload) ([]byte, error) {
	headAction := uint32(random.NewStdPRNG().GetUint64())
	newPld := payload.NewPayload(
		joinHead(headAction, pPld.GetHead()).uint64(),
		pPld.GetBody(),
	)

	actionKey := newActionKey(pRecv, headAction)

	p.setAction(actionKey)
	defer p.delAction(actionKey)

	// internal logger
	if err := p.enqueuePayload(pCtx, cIsRequest, pRecv, newPld); err != nil {
		return nil, fmt.Errorf("fetch payload: %w", err)
	}

	resp, err := p.recvResponse(actionKey)
	if err != nil {
		return nil, fmt.Errorf("receive response from fetch: %w", err)
	}

	return resp, nil
}

func (p *sNode) enqueueMessage(pCtx context.Context, pMsg message.IMessage) error {
	for i := uint64(0); i <= p.fSettings.GetRetryEnqueue(); i++ {
		if err := p.fQueue.EnqueueMessage(pMsg); err == nil {
			return nil
		}
		select {
		case <-pCtx.Done():
			return pCtx.Err()
		case <-time.After(p.fQueue.GetSettings().GetDuration()):
			// next iter
		}
	}
	return errors.New("enqueue message as send")
}

func (p *sNode) recvResponse(pActionKey string) ([]byte, error) {
	action, ok := p.getAction(pActionKey)
	if !ok {
		return nil, errors.New("action undefined")
	}
	select {
	case result, opened := <-action:
		if !opened {
			return nil, errors.New("chan is closed")
		}
		return result, nil
	case <-time.After(p.fSettings.GetFetchTimeWait()):
		return nil, errors.New("recv time is over")
	}
}

func (p *sNode) handleWrapper() network.IHandlerF {
	return func(pCtx context.Context, _ network.INode, pConn conn.IConn, pMsg net_message.IMessage) error {
		logBuilder := anon_logger.NewLogBuilder(p.fSettings.GetServiceName())

		// enrich logger
		logBuilder.WithConn(pConn)

		// try store hash of message
		if ok, err := p.storeHashWithBroadcast(pCtx, logBuilder, pMsg); !ok {
			// internal logger
			if err != nil {
				return fmt.Errorf("store hash with broadcast: %w", err)
			}
			return nil
		}

		client := p.fQueue.GetClient()

		// load encrypted message
		msg, err := message.LoadMessage(client.GetSettings(), pMsg.GetPayload().GetBody())
		if err != nil {
			// problem from sender's side
			p.fLogger.PushWarn(logBuilder.WithType(anon_logger.CLogWarnMessageNull))
			return fmt.Errorf("load message: %w", err)
		}

		// try decrypt message
		sender, pld, err := client.DecryptMessage(msg)
		if err != nil {
			p.fLogger.PushInfo(logBuilder.WithType(anon_logger.CLogInfoUndecryptable))
			return nil
		}

		// enrich logger
		logBuilder.WithPubKey(sender)

		// check sender in f2f list
		if !p.fFriends.InPubKeys(sender) {
			p.fLogger.PushWarn(logBuilder.WithType(anon_logger.CLogWarnNotFriend))
			return nil
		}

		// get header of payload
		head := loadHead(pld.GetHead())
		body := pld.GetBody()

		switch {
		// got response message from our side request
		case isResponse(body):
			// get session by payload head
			actionKey := newActionKey(sender, head.getAction())
			action, ok := p.getAction(actionKey)
			if !ok {
				p.fLogger.PushWarn(logBuilder.WithType(anon_logger.CLogBaseGetResponse))
				return nil
			}

			p.fLogger.PushInfo(logBuilder.WithType(anon_logger.CLogBaseGetResponse))
			action <- unwrapBytes(body)
			return nil

		// got request from another side (need generate response)
		case isRequest(body):
			// get function by payload head
			f, ok := p.getRoute(head.getRoute())
			if !ok || f == nil {
				p.fLogger.PushWarn(logBuilder.WithType(anon_logger.CLogWarnUnknownRoute))
				return nil
			}

			// response can be nil
			resp, err := f(pCtx, p, sender, unwrapBytes(body))
			if err != nil {
				p.fLogger.PushWarn(logBuilder.WithType(anon_logger.CLogWarnIncorrectResponse))
				return nil
			}
			if resp == nil {
				p.fLogger.PushInfo(logBuilder.WithType(anon_logger.CLogInfoWithoutResponse))
				return nil
			}

			// create the message and put this to the queue
			_ = p.enqueuePayload(pCtx, cIsResponse, sender, payload.NewPayload(pld.GetHead(), resp))
			// internal logger
			return nil

		// undefined type of message (not request/response)
		default:
			p.fLogger.PushWarn(logBuilder.WithType(anon_logger.CLogBaseMessageType))
			return nil
		}
	}
}

func (p *sNode) enqueuePayload(pCtx context.Context, pType iDataType, pRecv asymmetric.IPubKey, pPld payload.IPayload) error {
	logBuilder := anon_logger.NewLogBuilder(p.fSettings.GetServiceName())

	// enrich logger
	logBuilder.WithPubKey(p.fQueue.GetClient().GetPubKey())

	var (
		newBody []byte
		logType anon_logger.ILogType
	)

	switch pType {
	case cIsRequest:
		newBody = wrapRequest(pPld.GetBody())
		logType = anon_logger.CLogBaseEnqueueRequest
	case cIsResponse:
		newBody = wrapResponse(pPld.GetBody())
		logType = anon_logger.CLogBaseEnqueueResponse
	default:
		p.fLogger.PushErro(logBuilder.WithType(anon_logger.CLogBaseMessageType))
		return errors.New("unknown format type")
	}

	newPld := payload.NewPayload(pPld.GetHead(), newBody)
	msg, err := p.fQueue.GetClient().EncryptPayload(pRecv, newPld)
	if err != nil {
		p.fLogger.PushErro(logBuilder.WithType(anon_logger.CLogErroEncryptPayload))
		return fmt.Errorf("encrypt payload: %w", err)
	}

	var (
		size = len(msg.ToBytes())
		hash = msg.GetHash()
	)

	// enrich logger
	logBuilder.
		WithHash(hash).
		WithSize(size).
		WithType(logType)

	if err := p.enqueueMessage(pCtx, msg); err != nil {
		p.fLogger.PushErro(logBuilder)
		return fmt.Errorf("enqueue message: %w", err)
	}

	p.fLogger.PushInfo(logBuilder)
	return nil
}

func (p *sNode) storeHashWithBroadcast(pCtx context.Context, pLogBuilder anon_logger.ILogBuilder, pNetMsg net_message.IMessage) (bool, error) {
	var (
		size      = len(pNetMsg.GetPayload().GetBody())
		hash      = pNetMsg.GetHash()
		proof     = pNetMsg.GetProof()
		database  = p.fWrapperDB.Get()
		myAddress = p.fQueue.GetClient().GetPubKey().GetAddress().ToBytes()
	)

	// enrich logger
	pLogBuilder.
		WithHash(hash).
		WithProof(proof).
		WithSize(size)

	if database == nil {
		p.fLogger.PushErro(pLogBuilder.WithType(anon_logger.CLogErroDatabaseGet))
		return false, errors.New("database is nil")
	}

	// check already received data by hash
	allAddresses, err := database.Get(hash)
	hashIsExist := (err == nil)
	if hashIsExist && bytes.Contains(allAddresses, myAddress) {
		p.fLogger.PushInfo(pLogBuilder.WithType(anon_logger.CLogInfoExist))
		return false, nil
	}

	// set hash to database
	updateAddresses := bytes.Join(
		[][]byte{
			allAddresses,
			myAddress,
		},
		[]byte{},
	)
	if err := database.Set(hash, updateAddresses); err != nil {
		p.fLogger.PushErro(pLogBuilder.WithType(anon_logger.CLogErroDatabaseSet))
		return false, fmt.Errorf("database set: %w", err)
	}

	// do not send data if than already received
	if !hashIsExist {
		// broadcast message to network
		if err := p.networkBroadcast(pCtx, pNetMsg); err != nil {
			p.fLogger.PushWarn(pLogBuilder.WithType(anon_logger.CLogBaseBroadcast))
			// need pass error (some of connections may be closed)
		}
	}

	return true, nil
}

func (p *sNode) networkBroadcast(pCtx context.Context, pMsg net_message.IMessage) error {
	// redirect message to another nodes
	if err := p.fNetwork.BroadcastMessage(pCtx, pMsg); err != nil {
		return fmt.Errorf("network broadcast message: %w", err)
	}
	return nil
}

func (p *sNode) setRoute(pHead uint32, pHandle IHandlerF) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	p.fHandleRoutes[pHead] = pHandle
}

func (p *sNode) getRoute(pHead uint32) (IHandlerF, bool) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	f, ok := p.fHandleRoutes[pHead]
	return f, ok
}

func (p *sNode) getAction(pActionKey string) (chan []byte, bool) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	f, ok := p.fHandleActions[pActionKey]
	return f, ok
}

func (p *sNode) setAction(pActionKey string) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	p.fHandleActions[pActionKey] = make(chan []byte)
}

func (p *sNode) delAction(pActionKey string) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	delete(p.fHandleActions, pActionKey)
}

func newActionKey(pPubKey asymmetric.IPubKey, pHead uint32) string {
	pubKeyAddr := pPubKey.GetAddress().ToString()
	return fmt.Sprintf("%s-%d", pubKeyAddr, pHead)
}
