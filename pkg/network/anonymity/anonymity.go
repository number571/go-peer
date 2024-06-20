package anonymity

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/database"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/anonymity/queue"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/state"
	"github.com/number571/go-peer/pkg/utils"

	anon_logger "github.com/number571/go-peer/pkg/network/anonymity/logger"
	net_message "github.com/number571/go-peer/pkg/network/message"
)

var (
	_ INode = &sNode{}
)

type sNode struct {
	fMutex         sync.RWMutex
	fState         state.IState
	fSettings      ISettings
	fLogger        logger.ILogger
	fKVDatavase    database.IKVDatabase
	fNetwork       network.INode
	fQueue         queue.IMessageQueue
	fFriends       asymmetric.IListPubKeys
	fHandleRoutes  map[uint32]IHandlerF
	fHandleActions map[string]chan []byte
}

func NewNode(
	pSett ISettings,
	pLogger logger.ILogger,
	pKVDatavase database.IKVDatabase,
	pNetwork network.INode,
	pQueue queue.IMessageQueue,
	pFriends asymmetric.IListPubKeys,
) INode {
	return &sNode{
		fState:         state.NewBoolState(),
		fSettings:      pSett,
		fLogger:        pLogger,
		fKVDatavase:    pKVDatavase,
		fNetwork:       pNetwork,
		fQueue:         pQueue,
		fFriends:       pFriends,
		fHandleRoutes:  make(map[uint32]IHandlerF, 64),
		fHandleActions: make(map[string]chan []byte, 64),
	}
}

func (p *sNode) Run(pCtx context.Context) error {
	enableFunc := func() error {
		p.fNetwork.HandleFunc(
			p.fSettings.GetNetworkMask(),
			p.networkHandler,
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
			netMsg := p.fQueue.DequeueMessage(pCtx)
			if netMsg == nil {
				// context done
				break
			}

			logBuilder := anon_logger.NewLogBuilder(p.fSettings.GetServiceName())

			// update logger state
			p.enrichLogger(logBuilder, netMsg).
				WithPubKey(p.fQueue.GetClient().GetPubKey())

			// internal logger
			_, _ = p.storeHashWithBroadcast(pCtx, logBuilder, netMsg)
		}
	}
}

func (p *sNode) GetLogger() logger.ILogger {
	return p.fLogger
}

func (p *sNode) GetSettings() ISettings {
	return p.fSettings
}

func (p *sNode) GetKVDatabase() database.IKVDatabase {
	return p.fKVDatavase
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
func (p *sNode) SendPayload(
	pCtx context.Context,
	pRecv asymmetric.IPubKey,
	pPld payload.IPayload64,
) error {
	logBuilder := anon_logger.NewLogBuilder(p.fSettings.GetServiceName())
	if err := p.enqueuePayload(pCtx, logBuilder, pRecv, pPld); err != nil {
		// internal logger
		return utils.MergeErrors(ErrEnqueuePayload, err)
	}
	return nil
}

// Send message with response waiting.
// Payload head must be uint32.
func (p *sNode) FetchPayload(
	pCtx context.Context,
	pRecv asymmetric.IPubKey,
	pPld payload.IPayload32,
) ([]byte, error) {
	headAction := sAction(random.NewCSPRNG().GetUint64())
	actionKey := newActionKey(pRecv, headAction)

	p.setAction(actionKey)
	defer p.delAction(actionKey)

	newPld := payload.NewPayload64(
		joinHead(headAction.setType(true), pPld.GetHead()).uint64(),
		pPld.GetBody(),
	)

	logBuilder := anon_logger.NewLogBuilder(p.fSettings.GetServiceName())
	if err := p.enqueuePayload(pCtx, logBuilder, pRecv, newPld); err != nil {
		// internal logger
		return nil, utils.MergeErrors(ErrEnqueuePayload, err)
	}

	resp, err := p.recvResponse(pCtx, actionKey)
	if err != nil {
		return nil, utils.MergeErrors(ErrFetchResponse, err)
	}

	return resp, nil
}

func (p *sNode) enqueueMessage(pCtx context.Context, pMsg []byte) error {
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
	return ErrRetryLimit
}

func (p *sNode) recvResponse(pCtx context.Context, pActionKey string) ([]byte, error) {
	action, ok := p.getAction(pActionKey)
	if !ok {
		return nil, ErrActionIsNotFound
	}
	select {
	case <-pCtx.Done():
		return nil, pCtx.Err()
	case result, opened := <-action:
		if !opened {
			return nil, ErrActionIsClosed
		}
		return result, nil
	case <-time.After(p.fSettings.GetFetchTimeout()):
		return nil, ErrActionTimeout
	}
}

func (p *sNode) networkHandler(
	pCtx context.Context,
	_ network.INode, // used as p.fNetwork
	pConn conn.IConn,
	pNetMsg net_message.IMessage,
) error {
	logBuilder := anon_logger.NewLogBuilder(p.fSettings.GetServiceName())

	// update logger state
	p.enrichLogger(logBuilder, pNetMsg).
		WithConn(pConn)

	client := p.fQueue.GetClient()
	encMsg := pNetMsg.GetPayload().GetBody()

	// load encrypted message without decryption try
	if _, err := message.LoadMessage(client.GetSettings(), encMsg); err != nil {
		// problem from sender's side
		p.fLogger.PushWarn(logBuilder.WithType(anon_logger.CLogWarnMessageNull))
		return utils.MergeErrors(ErrLoadMessage, err)
	}

	// try store hash of message
	if ok, err := p.storeHashWithBroadcast(pCtx, logBuilder, pNetMsg); !ok {
		// internal logger
		if err != nil {
			return utils.MergeErrors(ErrStoreHashWithBroadcast, err)
		}
		return nil
	}

	// try decrypt message
	sender, decMsg, err := client.DecryptMessage(encMsg)
	if err != nil {
		p.fLogger.PushInfo(logBuilder.WithType(anon_logger.CLogInfoUndecryptable))
		return nil
	}

	// enrich logger
	logBuilder.WithPubKey(sender)

	// check sender's public key in f2f list
	if !p.fFriends.InPubKeys(sender) {
		// if f2f=enabled:
		// ignore reading message from unknown public key
		if !p.fSettings.GetF2FDisabled() {
			p.fLogger.PushWarn(logBuilder.WithType(anon_logger.CLogWarnNotFriend))
			return nil
		}
		// if f2f=disabled:
		// continue to read a message from unknown public key
		p.fLogger.PushInfo(logBuilder.WithType(anon_logger.CLogInfoPassF2FOption))
	}

	// get payload from decrypted message
	pld := payload.LoadPayload64(decMsg)
	if pld == nil {
		// got invalid payload64 format from sender
		p.fLogger.PushWarn(logBuilder.WithType(anon_logger.CLogWarnPayloadNull))
		return nil
	}

	// do request or response action
	return p.handleDoAction(pCtx, logBuilder, sender, pld)
}

func (p *sNode) handleDoAction(
	pCtx context.Context,
	pLogBuilder anon_logger.ILogBuilder,
	pSender asymmetric.IPubKey,
	pPld payload.IPayload64,
) error {
	// get [head:body] from payload
	head := loadHead(pPld.GetHead())
	body := pPld.GetBody()

	// check state of payload = [request,response]?
	action := head.getAction()

	if action.isRequest() {
		// got request from another side (need generate response)
		p.handleRequest(pCtx, pLogBuilder, pSender, head, body)
		return nil
	}

	// got response message from our side request
	p.handleResponse(pCtx, pLogBuilder, pSender, action, body)
	return nil
}

func (p *sNode) handleResponse(
	_ context.Context,
	pLogBuilder anon_logger.ILogBuilder,
	pSender asymmetric.IPubKey,
	pAction iAction,
	pBody []byte,
) {
	// get session by payload head
	actionKey := newActionKey(pSender, pAction)
	action, ok := p.getAction(actionKey)
	if !ok {
		p.fLogger.PushWarn(pLogBuilder.WithType(anon_logger.CLogBaseGetResponse))
		return
	}

	p.fLogger.PushInfo(pLogBuilder.WithType(anon_logger.CLogBaseGetResponse))
	action <- pBody
}

func (p *sNode) handleRequest(
	pCtx context.Context,
	pLogBuilder anon_logger.ILogBuilder,
	pSender asymmetric.IPubKey,
	pHead iHead,
	pBody []byte,
) {
	// get function by payload head
	f, ok := p.getRoute(pHead.getRoute())
	if !ok || f == nil {
		p.fLogger.PushWarn(pLogBuilder.WithType(anon_logger.CLogWarnUnknownRoute))
		return
	}

	// response can be nil
	resp, err := f(pCtx, p, pSender, pBody)
	if err != nil {
		p.fLogger.PushWarn(pLogBuilder.WithType(anon_logger.CLogWarnIncorrectResponse))
		return
	}
	if resp == nil {
		p.fLogger.PushInfo(pLogBuilder.WithType(anon_logger.CLogInfoWithoutResponse))
		return
	}

	// create response and put this to the queue
	// internal logger
	newHead := joinHead(pHead.getAction().setType(false), pHead.getRoute()).uint64()
	_ = p.enqueuePayload(
		pCtx,
		pLogBuilder,
		pSender,
		payload.NewPayload64(newHead, resp),
	)
}

func (p *sNode) enqueuePayload(
	pCtx context.Context,
	pLogBuilder anon_logger.ILogBuilder,
	pRecv asymmetric.IPubKey,
	pPld payload.IPayload64,
) error {
	client := p.fQueue.GetClient()
	isRequest := loadHead(pPld.GetHead()).getAction().isRequest()

	logType := anon_logger.CLogBaseEnqueueResponse
	if isRequest {
		logType = anon_logger.CLogBaseEnqueueRequest
		// enrich logger with raw message
		pLogBuilder.WithPubKey(client.GetPubKey())
	}

	msg, err := client.EncryptMessage(pRecv, pPld.ToBytes())
	if err != nil {
		p.fLogger.PushErro(pLogBuilder.WithType(anon_logger.CLogErroEncryptPayload))
		return utils.MergeErrors(ErrEncryptPayload, err)
	}

	if isRequest {
		// enrich logger with raw message
		// without hash: log only from net_message!
		pLogBuilder.WithSize(len(msg))
	}

	if err := p.enqueueMessage(pCtx, msg); err != nil {
		p.fLogger.PushWarn(pLogBuilder.WithType(logType))
		return utils.MergeErrors(ErrEnqueueMessage, err)
	}

	p.fLogger.PushInfo(pLogBuilder.WithType(logType))
	return nil
}

func (p *sNode) enrichLogger(pLogBuilder anon_logger.ILogBuilder, pNetMsg net_message.IMessage) anon_logger.ILogBuilder {
	var (
		size  = len(pNetMsg.ToBytes())
		hash  = pNetMsg.GetHash()
		proof = pNetMsg.GetProof()
	)
	return pLogBuilder.
		WithProof(proof).
		WithHash(hash).
		WithSize(size)
}

func (p *sNode) storeHashWithBroadcast(
	pCtx context.Context,
	pLogBuilder anon_logger.ILogBuilder,
	pNetMsg net_message.IMessage,
) (bool, error) {
	// try push hash into database
	hashIsSaved, err := p.storeHashIntoDatabase(pLogBuilder, pNetMsg.GetHash())
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
	// check already received data by hash
	if _, err := p.fKVDatavase.Get(pHash); err == nil {
		p.fLogger.PushInfo(pLogBuilder.WithType(anon_logger.CLogInfoExist))
		return false, nil
	}

	// set hash to database with new address
	if err := p.fKVDatavase.Set(pHash, []byte{}); err != nil {
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
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	f, ok := p.fHandleRoutes[pHead]
	return f, ok
}

func (p *sNode) getAction(pActionKey string) (chan []byte, bool) {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

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

func newActionKey(pPubKey asymmetric.IPubKey, pAction iAction) string {
	pubKeyAddr := pPubKey.GetHasher().ToString()
	return fmt.Sprintf("%s-%d", pubKeyAddr, pAction.uint31())
}
