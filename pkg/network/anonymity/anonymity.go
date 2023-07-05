package anonymity

import (
	"fmt"
	"strings"
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
	logger := anon_logger.NewLogger(p.fSettings.GetServiceName())
	if err := p.runQueue(logger); err != nil {
		return errors.WrapError(err, "run node")
	}
	p.fNetwork.HandleFunc(
		p.fSettings.GetNetworkMask(),
		p.handleWrapper(logger),
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

func (p *sNode) HandleMessage(pMsg message.IMessage) {
	logger := anon_logger.NewLogger(p.fSettings.GetServiceName())
	p.handleWrapper(logger)(
		p.fNetwork,
		nil,
		pMsg.ToBytes(),
	)
}

func (p *sNode) HandleFunc(pHead uint32, pHandle IHandlerF) INode {
	p.setRoute(pHead, pHandle)
	return p
}

// Send message without response waiting.
func (p *sNode) BroadcastPayload(pRecv asymmetric.IPubKey, pPld adapters.IPayload) error {
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

	if err := p.enqueuePayload(cIsRequest, pRecv, newPld); err != nil {
		return nil, errors.WrapError(err, "fetch payload")
	}

	resp, err := p.recv(actionKey)
	if err != nil {
		return nil, errors.WrapError(err, "receive response from fetch")
	}
	return resp, nil
}

func (p *sNode) enqueuePayload(pType iDataType, pRecv asymmetric.IPubKey, pPld payload.IPayload) error {
	if len(p.fNetwork.GetConnections()) == 0 {
		return errors.NewError("length of connections = 0")
	}

	var newBody []byte
	switch pType {
	case cIsRequest:
		newBody = wrapRequest(pPld.GetBody())
	case cIsResponse:
		newBody = wrapResponse(pPld.GetBody())
	default:
		return errors.NewError("unknown format type")
	}

	newPld := payload.NewPayload(pPld.GetHead(), newBody)
	msg, err := p.fQueue.GetClient().EncryptPayload(pRecv, newPld)
	if err != nil {
		return errors.WrapError(err, "encrypt payload")
	}

	if err := p.send(msg); err != nil {
		return errors.WrapError(err, "send message")
	}
	return nil
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

func (p *sNode) runQueue(pLogger anon_logger.ILogger) error {
	if err := p.fQueue.Run(); err != nil {
		return errors.WrapError(err, "run queue")
	}

	go func() {
		for {
			msg, ok := <-p.fQueue.DequeueMessage()
			if !ok {
				break
			}

			// store hash and push message to network
			if ok := p.storeHashWithBroadcast(pLogger, nil, msg); !ok {
				// internal logger
				continue
			}

			p.fLogger.PushInfo(pLogger.GetFmtLog(
				anon_logger.CLogBaseBroadcast,
				msg.GetBody().GetHash(),
				msg.GetBody().GetProof(),
				p.fQueue.GetClient().GetPubKey(),
				nil,
			))
		}
	}()

	return nil
}

func (p *sNode) handleWrapper(pLogger anon_logger.ILogger) network.IHandlerF {
	return func(_ network.INode, pConn conn.IConn, pMsgBytes []byte) {
		client := p.fQueue.GetClient()
		settings := client.GetSettings()

		msg := message.LoadMessage(
			message.NewSettings(&message.SSettings{
				FWorkSize:    settings.GetWorkSize(),
				FMessageSize: settings.GetMessageSize(),
			}),
			pMsgBytes,
		)

		// try store hash of message
		if ok := p.storeHashWithBroadcast(pLogger, pConn, msg); !ok {
			// internal logger
			return
		}

		var (
			hash  = msg.GetBody().GetHash()
			proof = msg.GetBody().GetProof()
		)

		// try decrypt message
		sender, pld, err := client.DecryptMessage(msg)
		if err != nil {
			p.fLogger.PushInfo(pLogger.GetFmtLog(anon_logger.CLogInfoUndecryptable, hash, proof, nil, pConn))
			return
		}

		// check sender in f2f list
		if !p.fFriends.InPubKeys(sender) {
			p.fLogger.PushWarn(pLogger.GetFmtLog(anon_logger.CLogWarnNotFriend, hash, proof, sender, pConn))
			return
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
				p.fLogger.PushWarn(pLogger.GetFmtLog(anon_logger.CLogInfoAction, hash, proof, sender, pConn))
				return
			}

			p.fLogger.PushInfo(pLogger.GetFmtLog(anon_logger.CLogInfoAction, hash, proof, sender, pConn))
			action <- unwrapBytes(body)
			return
		// got request from another side (need generate response)
		case isRequest(body):
			// get function by payload head
			f, ok := p.getRoute(head.getRoute())
			if !ok || f == nil {
				p.fLogger.PushWarn(pLogger.GetFmtLog(anon_logger.CLogWarnUnknownRoute, hash, proof, sender, pConn))
				return
			}

			// response can be nil
			resp, err := f(p, sender, hash, unwrapBytes(body))
			if err != nil {
				p.fLogger.PushWarn(pLogger.GetFmtLog(anon_logger.CLogWarnIncorrectResponse, hash, proof, sender, pConn))
				return
			}
			if resp == nil {
				p.fLogger.PushInfo(pLogger.GetFmtLog(anon_logger.CLogInfoWithoutResponse, hash, proof, sender, pConn))
				return
			}

			// create the message and put this to the queue
			if err := p.enqueuePayload(cIsResponse, sender, payload.NewPayload(pld.GetHead(), resp)); err != nil {
				p.fLogger.PushErro(pLogger.GetFmtLog(anon_logger.CLogBaseEnqueueResponse, hash, proof, sender, pConn))
				return
			}

			p.fLogger.PushInfo(pLogger.GetFmtLog(anon_logger.CLogBaseEnqueueResponse, hash, proof, sender, pConn))
			return
		// undefined type of message (not request/response)
		default:
			p.fLogger.PushErro(pLogger.GetFmtLog(anon_logger.CLogErroMessageType, hash, proof, sender, pConn))
			return
		}
	}
}

func (p *sNode) storeHashWithBroadcast(pLogger anon_logger.ILogger, pConn conn.IConn, pMsg message.IMessage) bool {
	if pMsg == nil {
		p.fLogger.PushWarn(pLogger.GetFmtLog(anon_logger.CLogWarnMessageNull, nil, 0, nil, pConn))
		return false
	}

	var (
		hash      = pMsg.GetBody().GetHash()
		proof     = pMsg.GetBody().GetProof()
		database  = p.fWrapperDB.Get()
		myAddress = p.fQueue.GetClient().GetPubKey().GetAddress().ToString()
	)

	if database == nil {
		p.fLogger.PushErro(pLogger.GetFmtLog(anon_logger.CLogErroDatabaseGet, hash, proof, nil, pConn))
		return false
	}

	hashDB := []byte(fmt.Sprintf("_hash_%X", hash))
	gotAddrs, err := database.Get(hashDB)

	// check already received data by hash
	hashIsExist := (err == nil)
	if hashIsExist && strings.Contains(string(gotAddrs), myAddress) {
		p.fLogger.PushInfo(pLogger.GetFmtLog(anon_logger.CLogInfoExist, hash, proof, nil, pConn))
		return false
	}

	// set hash to database
	updateAddrs := fmt.Sprintf("%s;%s", string(gotAddrs), myAddress)
	if err := database.Set(hashDB, []byte(updateAddrs)); err != nil {
		p.fLogger.PushErro(pLogger.GetFmtLog(anon_logger.CLogErroDatabaseSet, hash, proof, nil, pConn))
		return false
	}

	// do not send data if than already received
	if !hashIsExist {
		// broadcast message to network
		if err := p.networkBroadcast(pLogger, pMsg); err != nil {
			p.fLogger.PushWarn(pLogger.GetFmtLog(anon_logger.CLogBaseBroadcast, hash, proof, nil, nil))
			// need pass error (some of connections may be closed)
		}
	}

	return true
}

func (p *sNode) networkBroadcast(pLogger anon_logger.ILogger, pMsg message.IMessage) error {
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
	headString := fmt.Sprintf("%d", pHead)
	return fmt.Sprintf("%s-%s", pubKeyAddr, headString)
}
