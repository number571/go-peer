package anonymity

import (
	"context"
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
	"github.com/number571/go-peer/pkg/utils"

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
	fDBWrapper     IDBWrapper
	fNetwork       network.INode
	fQueue         queue.IMessageQueue
	fFriends       asymmetric.IListPubKeys
	fHandleRoutes  map[uint32]IHandlerF
	fHandleActions map[string]chan []byte
}

func NewNode(
	pSett ISettings,
	pLogger logger.ILogger,
	pDBWrapper IDBWrapper,
	pNetwork network.INode,
	pQueue queue.IMessageQueue,
	pFriends asymmetric.IListPubKeys,
) INode {
	return &sNode{
		fState:         state.NewBoolState(),
		fSettings:      pSett,
		fLogger:        pLogger,
		fDBWrapper:     pDBWrapper,
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
		return utils.MergeErrors(ErrRunning, err)
	}

	defer func() {
		disableFunc := func() error {
			p.fNetwork.HandleFunc(p.fSettings.GetNetworkMask(), nil)
			return nil
		}
		_ = p.fState.Disable(disableFunc)
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
				// context done
				break
			}

			logBuilder := anon_logger.NewLogBuilder(p.fSettings.GetServiceName())

			// enrich logger
			logBuilder.
				WithPubKey(p.fQueue.GetClient().GetPubKey())

			// internal logger
			_, _ = p.storeHashWithBroadcast(pCtx, logBuilder, msg)
		}
	}
}

func (p *sNode) GetLogger() logger.ILogger {
	return p.fLogger
}

func (p *sNode) GetSettings() ISettings {
	return p.fSettings
}

func (p *sNode) GetDBWrapper() IDBWrapper {
	return p.fDBWrapper
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
func (p *sNode) SendPayload(pCtx context.Context, pRecv asymmetric.IPubKey, pPld payload.IPayload) error {
	if err := p.enqueuePayload(pCtx, cIsRequest, pRecv, pPld); err != nil {
		// internal logger
		return utils.MergeErrors(ErrBroadcastPayload, err)
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

	if err := p.enqueuePayload(pCtx, cIsRequest, pRecv, newPld); err != nil {
		// internal logger
		return nil, utils.MergeErrors(ErrEnqueuePayload, err)
	}

	resp, err := p.recvResponse(actionKey)
	if err != nil {
		return nil, utils.MergeErrors(ErrFetchResponse, err)
	}

	return resp, nil
}

func (p *sNode) enqueueMessage(pCtx context.Context, pMsg message.IMessage) error {
	retryNum := p.fSettings.GetRetryEnqueue()
	for i := uint64(0); i <= retryNum; i++ {
		if err := p.fQueue.EnqueueMessage(pMsg); err == nil {
			return nil
		}
		if i == retryNum {
			break
		}
		select {
		case <-pCtx.Done():
			return pCtx.Err()
		case <-time.After(p.fQueue.GetSettings().GetDuration()):
			// next iter
		}
	}
	return ErrEnqueueMessage
}

func (p *sNode) recvResponse(pActionKey string) ([]byte, error) {
	action, ok := p.getAction(pActionKey)
	if !ok {
		return nil, ErrActionIsNotFound
	}
	select {
	case result, opened := <-action:
		if !opened {
			return nil, ErrActionIsClosed
		}
		return result, nil
	case <-time.After(p.fSettings.GetFetchTimeWait()):
		return nil, ErrActionTimeout
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
				return utils.MergeErrors(ErrStoreHashWithBroadcast, err)
			}
			return nil
		}

		client := p.fQueue.GetClient()

		// load encrypted message
		msg, err := message.LoadMessage(client.GetSettings(), pMsg.GetPayload().GetBody())
		if err != nil {
			// problem from sender's side
			p.fLogger.PushWarn(logBuilder.WithType(anon_logger.CLogWarnMessageNull))
			return utils.MergeErrors(ErrLoadMessage, err)
		}

		// try decrypt message
		sender, pld, err := client.DecryptMessage(msg)
		if err != nil {
			p.fLogger.PushInfo(logBuilder.WithType(anon_logger.CLogInfoUndecryptable))
			return nil
		}

		// enrich logger
		logBuilder.WithPubKey(sender)

		// check sender's public key in f2f list
		if !p.fFriends.InPubKeys(sender) {
			switch p.fSettings.GetF2FDisabled() {
			case true:
				// continue to read a message from unknown public key
				p.fLogger.PushInfo(logBuilder.WithType(anon_logger.CLogInfoPassF2FOption))
			default:
				// ignore reading messages from unknown public key
				p.fLogger.PushWarn(logBuilder.WithType(anon_logger.CLogWarnNotFriend))
				return nil
			}
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

			// create response and put this to the queue
			// internal logger
			_ = p.enqueuePayload(pCtx, cIsResponse, sender, payload.NewPayload(pld.GetHead(), resp))
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
		return ErrUnknownType
	}

	newPld := payload.NewPayload(pPld.GetHead(), newBody)
	msg, err := p.fQueue.GetClient().EncryptPayload(pRecv, newPld)
	if err != nil {
		p.fLogger.PushErro(logBuilder.WithType(anon_logger.CLogErroEncryptPayload))
		return utils.MergeErrors(ErrEncryptPayload, err)
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
		return err
	}

	p.fLogger.PushInfo(logBuilder)
	return nil
}

func (p *sNode) storeHashWithBroadcast(pCtx context.Context, pLogBuilder anon_logger.ILogBuilder, pNetMsg net_message.IMessage) (bool, error) {
	var (
		size  = len(pNetMsg.GetPayload().GetBody())
		hash  = pNetMsg.GetHash()
		proof = pNetMsg.GetProof()
	)

	// enrich logger
	pLogBuilder.
		WithHash(hash).
		WithProof(proof).
		WithSize(size)

	// try push hash into database
	hashIsSaved, err := p.storeHashIntoDatabase(pLogBuilder, hash)
	if err != nil || !hashIsSaved {
		// internal logger
		return false, err
	}

	// redirect message to another nodes
	if err := p.fNetwork.BroadcastMessage(pCtx, pNetMsg); err != nil {
		// some connections can return errors
		p.fLogger.PushWarn(pLogBuilder.WithType(anon_logger.CLogBaseBroadcast))
		return true, nil
	}

	// full success broadcast
	p.fLogger.PushInfo(pLogBuilder.WithType(anon_logger.CLogBaseBroadcast))
	return true, nil
}

func (p *sNode) storeHashIntoDatabase(pLogBuilder anon_logger.ILogBuilder, pHash []byte) (bool, error) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	database := p.fDBWrapper.Get()
	if database == nil {
		p.fLogger.PushErro(pLogBuilder.WithType(anon_logger.CLogErroDatabaseGet))
		return false, ErrNilDB
	}

	// check already received data by hash
	if _, err := database.Get(pHash); err == nil {
		p.fLogger.PushInfo(pLogBuilder.WithType(anon_logger.CLogInfoExist))
		return false, nil
	}

	// set hash to database with new address
	if err := database.Set(pHash, []byte{}); err != nil {
		p.fLogger.PushErro(pLogBuilder.WithType(anon_logger.CLogErroDatabaseSet))
		return false, utils.MergeErrors(ErrSetHashIntoDB, err)
	}

	return true, nil
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
	pubKeyAddr := pPubKey.GetHasher().ToString()
	return fmt.Sprintf("%s-%d", pubKeyAddr, pHead)
}
