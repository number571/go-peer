package qb

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/number571/go-peer/pkg/anonymity/qb/adapters"
	"github.com/number571/go-peer/pkg/anonymity/qb/queue"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/message/layer1"
	"github.com/number571/go-peer/pkg/message/layer2"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/state"
	"github.com/number571/go-peer/pkg/storage/database"

	anon_logger "github.com/number571/go-peer/pkg/anonymity/qb/logger"
)

var (
	_ INode = &sNode{}
)

type sNode struct {
	fMutex         sync.RWMutex
	fState         state.IState
	fSettings      ISettings
	fLogger        logger.ILogger
	fAdapter       adapters.IAdapter
	fKVDatavase    database.IKVDatabase
	fQBProcessor   queue.IQBProblemProcessor
	fMapPubKeys    asymmetric.IMapPubKeys
	fHandleRoutes  map[uint32]IHandlerF
	fHandleActions map[string]chan []byte
}

func NewNode(
	pSett ISettings,
	pLogger logger.ILogger,
	pAdapter adapters.IAdapter,
	pKVDatavase database.IKVDatabase,
	pQBProcessor queue.IQBProblemProcessor,
) INode {
	return &sNode{
		fState:         state.NewBoolState(),
		fSettings:      pSett,
		fLogger:        pLogger,
		fAdapter:       pAdapter,
		fKVDatavase:    pKVDatavase,
		fQBProcessor:   pQBProcessor,
		fMapPubKeys:    asymmetric.NewMapPubKeys(),
		fHandleRoutes:  make(map[uint32]IHandlerF, 64),
		fHandleActions: make(map[string]chan []byte, 64),
	}
}

func (p *sNode) Run(pCtx context.Context) error {
	if err := p.fState.Enable(nil); err != nil {
		return errors.Join(ErrRunning, err)
	}
	defer func() { _ = p.fState.Disable(nil) }()

	chCtx, cancel := context.WithCancel(pCtx)
	defer cancel()

	const N = 3

	errs := make([]error, N)
	wg := &sync.WaitGroup{}
	wg.Add(N)

	go func() {
		defer func() { wg.Done(); cancel() }()
		errs[0] = p.fQBProcessor.Run(chCtx)
	}()
	go func() {
		defer func() { wg.Done(); cancel() }()
		errs[1] = p.runConsumer(chCtx)
	}()
	go func() {
		defer func() { wg.Done(); cancel() }()
		errs[2] = p.runProducer(chCtx)
	}()

	wg.Wait()

	select {
	case <-pCtx.Done():
		return pCtx.Err()
	default:
		return errors.Join(errs...)
	}
}

func (p *sNode) runProducer(pCtx context.Context) error {
	for {
		select {
		case <-pCtx.Done():
			return pCtx.Err()
		default:
			netMsg := p.fQBProcessor.DequeueMessage(pCtx)
			if netMsg == nil {
				// context done
				continue
			}
			// internal logger
			_, _ = p.produceMessage(pCtx, netMsg)
		}
	}
}

func (p *sNode) runConsumer(pCtx context.Context) error {
	for {
		select {
		case <-pCtx.Done():
			return pCtx.Err()
		default:
			netMsg, err := p.fAdapter.Consume(pCtx)
			if err != nil {
				// context done or error
				continue
			}
			// internal logger
			_ = p.consumeMessage(pCtx, netMsg)
		}
	}
}

func (p *sNode) GetLogger() logger.ILogger {
	return p.fLogger
}

func (p *sNode) GetSettings() ISettings {
	return p.fSettings
}

func (p *sNode) GetAdapter() adapters.IAdapter {
	return p.fAdapter
}

func (p *sNode) GetKVDatabase() database.IKVDatabase {
	return p.fKVDatavase
}

func (p *sNode) GetQBProcessor() queue.IQBProblemProcessor {
	return p.fQBProcessor
}

// Return f2f structure.
func (p *sNode) GetMapPubKeys() asymmetric.IMapPubKeys {
	return p.fMapPubKeys
}

func (p *sNode) HandleFunc(pHead uint32, pHandle IHandlerF) INode {
	p.setRoute(pHead, pHandle)
	return p
}

// Send message without response waiting.
func (p *sNode) SendPayload(
	_ context.Context,
	pRecv asymmetric.IPubKey,
	pPld payload.IPayload64,
) error {
	logBuilder := anon_logger.NewLogBuilder(p.fSettings.GetServiceName())
	if err := p.enqueuePayload(logBuilder, pRecv, pPld); err != nil {
		// internal logger
		return errors.Join(ErrEnqueuePayload, err)
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
	headAction := sAction(random.NewRandom().GetUint64()) //nolint:gosec
	actionKey := newActionKey(pRecv, headAction)

	p.setAction(actionKey)
	defer p.delAction(actionKey)

	newPld := payload.NewPayload64(
		joinHead(headAction.setType(true), pPld.GetHead()).uint64(),
		pPld.GetBody(),
	)

	logBuilder := anon_logger.NewLogBuilder(p.fSettings.GetServiceName())
	if err := p.enqueuePayload(logBuilder, pRecv, newPld); err != nil {
		// internal logger
		return nil, errors.Join(ErrEnqueuePayload, err)
	}

	resp, err := p.recvResponse(pCtx, actionKey)
	if err != nil {
		return nil, errors.Join(ErrFetchResponse, err)
	}

	return resp, nil
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

func (p *sNode) consumeMessage(pCtx context.Context, pNetMsg layer1.IMessage) error {
	logBuilder := anon_logger.NewLogBuilder(p.fSettings.GetServiceName())

	// update logger state
	p.enrichLogger(logBuilder, pNetMsg)

	// check network message on correct format
	if ok := p.checkMessageLayer1(pNetMsg); !ok {
		// another network mask || message settings
		p.fLogger.PushWarn(logBuilder.WithType(anon_logger.CLogWarnIncorrectLayer1))
		return ErrInvalidLayer1Message
	}

	client := p.fQBProcessor.GetClient()
	encMsg := pNetMsg.GetPayload().GetBody()

	// load encrypted message without decryption try
	if _, err := layer2.LoadMessage(client.GetMessageSize(), encMsg); err != nil {
		// problem from sender's side (invalid structure)
		p.fLogger.PushWarn(logBuilder.WithType(anon_logger.CLogWarnMessageNull))
		return errors.Join(ErrLoadMessage, err)
	}

	// try store hash of message
	if err := p.storeHashIntoDatabase(logBuilder, pNetMsg); err != nil {
		// internal logger
		if errors.Is(err, ErrHashAlreadyExist) {
			return nil
		}
		return errors.Join(ErrStoreHashIntoDatabase, err)
	}

	// try decrypt consumed message
	pubKey, decMsg, err := client.DecryptMessage(p.fMapPubKeys, encMsg)
	if err != nil {
		p.fLogger.PushInfo(logBuilder.WithType(anon_logger.CLogInfoUndecryptable))
		return nil
	}

	// enrich logger
	logBuilder.WithPubKey(pubKey)

	// get payload from decrypted message
	pld := payload.LoadPayload64(decMsg)
	if pld == nil {
		// got invalid payload64 format from sender
		p.fLogger.PushWarn(logBuilder.WithType(anon_logger.CLogWarnPayloadNull))
		return nil
	}

	// do request or response action
	return p.handleDoAction(pCtx, logBuilder, pubKey, pld)
}

func (p *sNode) checkMessageLayer1(pNetMsg layer1.IMessage) bool {
	settings := p.fQBProcessor.GetSettings()
	if settings.GetNetworkMask() != pNetMsg.GetPayload().GetHead() {
		return false
	}
	_, err := layer1.LoadMessage(
		settings.GetMessageConstructSettings().GetSettings(),
		pNetMsg.ToBytes(),
	)
	return err == nil
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
		pLogBuilder,
		pSender,
		payload.NewPayload64(newHead, resp),
	)
}

func (p *sNode) enqueuePayload(
	pLogBuilder anon_logger.ILogBuilder,
	pRecv asymmetric.IPubKey,
	pPld payload.IPayload64,
) error {
	logType := anon_logger.CLogBaseEnqueueResponse
	pldBytes := pPld.ToBytes()

	if loadHead(pPld.GetHead()).getAction().isRequest() {
		logType = anon_logger.CLogBaseEnqueueRequest
		client := p.fQBProcessor.GetClient()
		// enrich logger
		pLogBuilder.
			WithPubKey(client.GetPrivKey().GetPubKey()).
			WithSize(len(pldBytes))
	}

	if err := p.fQBProcessor.EnqueueMessage(pRecv, pldBytes); err != nil {
		p.fLogger.PushWarn(pLogBuilder.WithType(logType))
		return errors.Join(ErrEnqueueMessage, err)
	}

	p.fLogger.PushInfo(pLogBuilder.WithType(logType))
	return nil
}

func (p *sNode) enrichLogger(pLogBuilder anon_logger.ILogBuilder, pNetMsg layer1.IMessage) anon_logger.ILogBuilder {
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

func (p *sNode) produceMessage(
	pCtx context.Context,
	pNetMsg layer1.IMessage,
) (bool, error) {
	serviceName := p.fSettings.GetServiceName()

	// create logger state
	logBuilder := p.enrichLogger(anon_logger.NewLogBuilder(serviceName), pNetMsg).
		WithPubKey(p.fQBProcessor.GetClient().GetPrivKey().GetPubKey())

	// try push hash into database
	if err := p.storeHashIntoDatabase(logBuilder, pNetMsg); err != nil {
		// internal logger
		if errors.Is(err, ErrHashAlreadyExist) {
			return false, nil
		}
		return false, err
	}

	// redirect message to another nodes
	if err := p.fAdapter.Produce(pCtx, pNetMsg); err != nil {
		// some connections can return errors
		p.fLogger.PushWarn(logBuilder.WithType(anon_logger.CLogBaseBroadcast))
		return true, nil
	}

	// full success broadcast
	p.fLogger.PushInfo(logBuilder.WithType(anon_logger.CLogBaseBroadcast))
	return true, nil
}

func (p *sNode) storeHashIntoDatabase(pLogBuilder anon_logger.ILogBuilder, pNetMsg layer1.IMessage) error {
	// check already received data by hash
	hash := hashing.NewHasher(pNetMsg.GetPayload().GetBody()).ToBytes()
	_, err := p.fKVDatavase.Get(hash)
	if err == nil {
		p.fLogger.PushInfo(pLogBuilder.WithType(anon_logger.CLogInfoExist))
		return ErrHashAlreadyExist
	}
	if !errors.Is(err, database.ErrNotFound) {
		p.fLogger.PushInfo(pLogBuilder.WithType(anon_logger.CLogErroDatabaseGet))
		return errors.Join(ErrGetHashFromDB, err)
	}
	// set hash to database with new address
	if err := p.fKVDatavase.Set(hash, []byte{}); err != nil {
		p.fLogger.PushErro(pLogBuilder.WithType(anon_logger.CLogErroDatabaseSet))
		return errors.Join(ErrSetHashIntoDB, err)
	}
	return nil
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
	pubKeyAddr := hashing.NewHasher(pPubKey.ToBytes()).ToBytes()
	return fmt.Sprintf("%s-%d", pubKeyAddr, pAction.uint31())
}
