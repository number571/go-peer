package anonymity

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/client/queue"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/payload"

	"github.com/number571/go-peer/pkg/network/anonymity/adapters"
	anon_logger "github.com/number571/go-peer/pkg/network/anonymity/logger"
)

var (
	_ INode = &sNode{}
)

type sNode struct {
	fMutex         sync.Mutex
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

func (p *sNode) Run() error {
	if err := p.runQueue(); err != nil {
		return errors.WrapError(err, "run node")
	}
	p.fNetwork.HandleFunc(
		p.fSettings.GetNetworkMask(),
		p.handleWrapper(),
	)
	return nil
}

func (p *sNode) Stop() error {
	p.fNetwork.HandleFunc(p.fSettings.GetNetworkMask(), nil)
	if err := p.fQueue.Stop(); err != nil {
		return errors.WrapError(err, "stop node")
	}
	return nil
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

func (p *sNode) HandleMessage(pMsg message.IMessage) error {
	handler := p.handleWrapper()
	return handler(p.fNetwork, nil, pMsg.ToBytes())
}

func (p *sNode) HandleFunc(pHead uint32, pHandle IHandlerF) INode {
	p.setRoute(pHead, pHandle)
	return p
}

// Send message without response waiting.
func (p *sNode) BroadcastPayload(pRecv asymmetric.IPubKey, pPld adapters.IPayload) error {
	// internal logger
	if err := p.enqueuePayload(cIsRequest, pRecv, pPld.ToOrigin()); err != nil {
		return errors.WrapError(err, "broadcast payload")
	}
	return nil
}

// Send message with response waiting.
// Payload head must be uint32.
func (p *sNode) FetchPayload(pRecv asymmetric.IPubKey, pPld adapters.IPayload) ([]byte, error) {
	headAction := uint32(random.NewStdPRNG().GetUint64())
	newPld := payload.NewPayload(
		joinHead(headAction, pPld.GetHead()).uint64(),
		pPld.GetBody(),
	)

	actionKey := newActionKey(pRecv, headAction)

	p.setAction(actionKey)
	defer p.delAction(actionKey)

	// internal logger
	if err := p.enqueuePayload(cIsRequest, pRecv, newPld); err != nil {
		return nil, errors.WrapError(err, "fetch payload")
	}

	resp, err := p.recv(actionKey)
	if err != nil {
		return nil, errors.WrapError(err, "receive response from fetch")
	}

	return resp, nil
}

func (p *sNode) send(pMsg message.IMessage) error {
	for i := uint64(0); i <= p.fSettings.GetRetryEnqueue(); i++ {
		if err := p.fQueue.EnqueueMessage(pMsg); err != nil {
			time.Sleep(p.fQueue.GetSettings().GetDuration())
			continue
		}
		return nil
	}
	return errors.NewError("enqueue message as send")
}

func (p *sNode) recv(pActionKey string) ([]byte, error) {
	action, ok := p.getAction(pActionKey)
	if !ok {
		return nil, errors.NewError("action undefined")
	}
	select {
	case result, opened := <-action:
		if !opened {
			return nil, errors.NewError("chan is closed")
		}
		return result, nil
	case <-time.After(p.fSettings.GetFetchTimeWait()):
		return nil, errors.NewError("recv time is over")
	}
}

func (p *sNode) runQueue() error {
	if err := p.fQueue.Run(); err != nil {
		return errors.WrapError(err, "run queue")
	}

	go func() {
		for {
			msg, ok := <-p.fQueue.DequeueMessage()
			if !ok {
				break
			}

			logBuilder := anon_logger.NewLogBuilder(p.fSettings.GetServiceName())

			// store hash and push message to network
			if ok, _ := p.storeHashWithBroadcast(logBuilder, msg); !ok {
				// internal logger
				continue
			}

			// enrich logger
			logBuilder.
				WithPubKey(p.fQueue.GetClient().GetPubKey()).
				WithType(anon_logger.CLogBaseBroadcast)

			p.fLogger.PushInfo(logBuilder)
		}
	}()

	return nil
}

func (p *sNode) handleWrapper() network.IHandlerF {
	return func(_ network.INode, pConn conn.IConn, pMsgBytes []byte) error {
		logBuilder := anon_logger.NewLogBuilder(p.fSettings.GetServiceName())

		// enrich logger
		logBuilder.WithConn(pConn)

		client := p.fQueue.GetClient()
		settings := client.GetSettings()

		msg := message.LoadMessage(
			message.NewSettings(&message.SSettings{
				FWorkSizeBits:     settings.GetWorkSizeBits(),
				FMessageSizeBytes: settings.GetMessageSizeBytes(),
			}),
			pMsgBytes,
		)

		// try store hash of message
		if ok, err := p.storeHashWithBroadcast(logBuilder, msg); !ok {
			// internal logger
			if err != nil {
				return errors.WrapError(err, "store hash with broadcast")
			}
			return nil
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
			resp, err := f(p, sender, unwrapBytes(body))
			if err != nil {
				p.fLogger.PushWarn(logBuilder.WithType(anon_logger.CLogWarnIncorrectResponse))
				return nil
			}
			if resp == nil {
				p.fLogger.PushInfo(logBuilder.WithType(anon_logger.CLogInfoWithoutResponse))
				return nil
			}

			// create the message and put this to the queue
			// internal logger
			_ = p.enqueuePayload(cIsResponse, sender, payload.NewPayload(pld.GetHead(), resp))
			return nil

		// undefined type of message (not request/response)
		default:
			p.fLogger.PushWarn(logBuilder.WithType(anon_logger.CLogBaseMessageType))
			return nil
		}
	}
}

func (p *sNode) enqueuePayload(pType iDataType, pRecv asymmetric.IPubKey, pPld payload.IPayload) error {
	logBuilder := anon_logger.NewLogBuilder(p.fSettings.GetServiceName())

	// enrich logger
	logBuilder.WithPubKey(p.fQueue.GetClient().GetPubKey())

	if len(p.fNetwork.GetConnections()) == 0 {
		p.fLogger.PushWarn(logBuilder.WithType(anon_logger.CLogWarnNotConnection))
		return errors.NewError("length of connections = 0")
	}

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
		return errors.NewError("unknown format type")
	}

	newPld := payload.NewPayload(pPld.GetHead(), newBody)
	msg, err := p.fQueue.GetClient().EncryptPayload(pRecv, newPld)
	if err != nil {
		p.fLogger.PushErro(logBuilder.WithType(anon_logger.CLogErroEncryptPayload))
		return errors.WrapError(err, "encrypt payload")
	}

	var (
		size  = len(msg.ToBytes())
		hash  = msg.GetBody().GetHash()
		proof = msg.GetBody().GetProof()
	)

	// enrich logger
	logBuilder.
		WithHash(hash).
		WithProof(proof).
		WithSize(size).
		WithType(logType)

	if err := p.send(msg); err != nil {
		p.fLogger.PushErro(logBuilder)
		return errors.WrapError(err, "send message")
	}

	p.fLogger.PushInfo(logBuilder)
	return nil
}

func (p *sNode) storeHashWithBroadcast(pLogBuilder anon_logger.ILogBuilder, pMsg message.IMessage) (bool, error) {
	if pMsg == nil {
		p.fLogger.PushWarn(pLogBuilder.WithType(anon_logger.CLogWarnMessageNull))
		return false, errors.NewError("message is nil")
	}

	var (
		size      = len(pMsg.ToBytes())
		hash      = pMsg.GetBody().GetHash()
		proof     = pMsg.GetBody().GetProof()
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
		return false, errors.NewError("database is nil")
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
		return false, errors.WrapError(err, "database set")
	}

	// do not send data if than already received
	if !hashIsExist {
		// broadcast message to network
		if err := p.networkBroadcast(pMsg); err != nil {
			p.fLogger.PushWarn(pLogBuilder.WithType(anon_logger.CLogBaseBroadcast))
			// need pass error (some of connections may be closed)
		}
	}

	return true, nil
}

func (p *sNode) networkBroadcast(pMsg message.IMessage) error {
	// redirect message to another nodes
	err := p.fNetwork.BroadcastPayload(
		payload.NewPayload(
			p.fSettings.GetNetworkMask(),
			pMsg.ToBytes(),
		),
	)
	if err != nil {
		return errors.WrapError(err, "network broadcast payload")
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
