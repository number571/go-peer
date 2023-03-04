package anonymity

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/client/queue"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/go-peer/pkg/types"

	anon_logger "github.com/number571/go-peer/pkg/network/anonymity/logger"
)

var (
	_ INode = &sNode{}
)

type sNode struct {
	fMutex         sync.Mutex
	fSettings      ISettings
	fLogger        logger.ILogger
	fKeyValueDB    database.IKeyValueDB
	fNetwork       network.INode
	fQueue         queue.IMessageQueue
	fFriends       asymmetric.IListPubKeys
	fHandleRoutes  map[uint32]IHandlerF
	fHandleActions map[string]chan []byte
}

func NewNode(
	sett ISettings,
	log logger.ILogger,
	kvDB database.IKeyValueDB,
	nnode network.INode,
	queue queue.IMessageQueue,
	friends asymmetric.IListPubKeys,
) INode {
	return &sNode{
		fSettings:      sett,
		fLogger:        log,
		fKeyValueDB:    kvDB,
		fNetwork:       nnode,
		fQueue:         queue,
		fFriends:       friends,
		fHandleRoutes:  make(map[uint32]IHandlerF),
		fHandleActions: make(map[string]chan []byte),
	}
}

func (node *sNode) Run() error {
	logger := anon_logger.NewLogger(node.GetSettings().GetServiceName())

	if err := node.runQueue(logger); err != nil {
		return err
	}

	node.GetNetworkNode().HandleFunc(
		node.GetSettings().GetNetworkMask(),
		node.handleWrapper(logger),
	)

	return nil
}

func (node *sNode) Stop() error {
	node.GetNetworkNode().HandleFunc(node.GetSettings().GetNetworkMask(), nil)
	return types.StopAllCommands([]types.ICommand{
		node.GetMessageQueue(),
	})
}

func (node *sNode) GetLogger() logger.ILogger {
	return node.fLogger
}

func (node *sNode) GetSettings() ISettings {
	return node.fSettings
}

func (node *sNode) GetKeyValueDB() database.IKeyValueDB {
	return node.fKeyValueDB
}

func (node *sNode) GetNetworkNode() network.INode {
	return node.fNetwork
}

func (node *sNode) GetMessageQueue() queue.IMessageQueue {
	return node.fQueue
}

// Return f2f structure.
func (node *sNode) GetListPubKeys() asymmetric.IListPubKeys {
	return node.fFriends
}

func (node *sNode) HandleFunc(head uint32, handle IHandlerF) INode {
	node.setRoute(head, handle)
	return node
}

// Send message without response waiting.
func (node *sNode) BroadcastPayload(recv asymmetric.IPubKey, pld payload.IPayload) error {
	if len(node.GetNetworkNode().GetConnections()) == 0 {
		return errors.New("length of connections = 0")
	}

	msg, err := node.GetMessageQueue().GetClient().EncryptPayload(recv, pld)
	if err != nil {
		return err
	}

	return node.send(msg)
}

// Send message with response waiting.
// Payload head must be uint32.
func (node *sNode) FetchPayload(recv asymmetric.IPubKey, pld payload.IPayload) ([]byte, error) {
	headAction := uint32(random.NewStdPRNG().GetUint64())
	headRoute := mustBeUint32(pld.GetHead())

	newPld := payload.NewPayload(
		joinHead(headAction, headRoute).Uint64(),
		pld.GetBody(),
	)

	actionKey := newActionKey(recv, headAction)

	node.setAction(actionKey)
	defer node.delAction(actionKey)

	if err := node.BroadcastPayload(recv, newPld); err != nil {
		return nil, err
	}
	return node.recv(actionKey)
}

func (node *sNode) send(msg message.IMessage) error {
	for i := uint64(0); i <= node.GetSettings().GetRetryEnqueue(); i++ {
		if err := node.GetMessageQueue().EnqueueMessage(msg); err != nil {
			time.Sleep(node.GetMessageQueue().GetSettings().GetDuration())
			continue
		}
		return nil
	}
	return fmt.Errorf("failed: enqueue message")
}

func (node *sNode) recv(actionKey string) ([]byte, error) {
	action, ok := node.getAction(actionKey)
	if !ok {
		return nil, errors.New("action undefined")
	}
	select {
	case result, opened := <-action:
		if !opened {
			return nil, errors.New("chan is closed")
		}
		return result, nil
	case <-time.After(node.GetSettings().GetTimeWait()):
		return nil, errors.New("time is over")
	}
}

func (node *sNode) runQueue(logger anon_logger.ILogger) error {
	if err := node.GetMessageQueue().Run(); err != nil {
		return err
	}

	go func() {
		for {
			msg, ok := <-node.GetMessageQueue().DequeueMessage()
			if !ok {
				break
			}

			var (
				hash   = msg.GetBody().GetHash()
				proof  = msg.GetBody().GetProof()
				pubKey = node.GetMessageQueue().GetClient().GetPubKey()
			)

			if err := node.networkBroadcast(logger, msg); err != nil {
				node.fLogger.PushErro(logger.GetFmtLog(anon_logger.CLogErroMiddleware, hash, proof, pubKey, nil))
				continue
			}

			node.fLogger.PushInfo(logger.GetFmtLog(anon_logger.CLogBaseBroadcast, hash, proof, pubKey, nil))
		}
	}()

	return nil
}

func (node *sNode) handleWrapper(logger anon_logger.ILogger) network.IHandlerF {
	return func(_ network.INode, conn conn.IConn, reqBytes []byte) {
		msg := message.LoadMessage(
			reqBytes,
			message.NewParams(
				node.GetMessageQueue().GetClient().GetSettings().GetMessageSize(),
				node.GetMessageQueue().GetClient().GetSettings().GetWorkSize(),
			),
		)
		if msg == nil {
			node.GetLogger().PushWarn(logger.GetFmtLog(anon_logger.CLogWarnMessageNull, nil, 0, nil, conn))
			return
		}

		var (
			addr  = node.GetMessageQueue().GetClient().GetPubKey().Address().ToString()
			hash  = msg.GetBody().GetHash()
			proof = msg.GetBody().GetProof()
		)

		hashDB := []byte(fmt.Sprintf("_hash_%X", hash))
		gotAddrs, err := node.GetKeyValueDB().Get(hashDB)

		// check already received data by hash
		hashIsExist := (err == nil)
		if hashIsExist && strings.Contains(string(gotAddrs), addr) {
			node.GetLogger().PushInfo(logger.GetFmtLog(anon_logger.CLogInfoExist, hash, proof, nil, conn))
			return
		}

		// set hash to database
		updateAddrs := fmt.Sprintf("%s;%s", string(gotAddrs), addr)
		if err := node.GetKeyValueDB().Set(hashDB, []byte(updateAddrs)); err != nil {
			node.GetLogger().PushErro(logger.GetFmtLog(anon_logger.CLogErroDatabaseSet, hash, proof, nil, conn))
			return
		}

		// do not send data if than already received
		if !hashIsExist {
			// broadcast message to network
			if err := node.networkBroadcast(logger, msg); err != nil {
				node.GetLogger().PushErro(logger.GetFmtLog(anon_logger.CLogErroMiddleware, hash, proof, nil, conn))
				return
			}
		}

		// try decrypt message
		sender, pld, err := node.GetMessageQueue().GetClient().DecryptMessage(msg)
		if err != nil {
			node.GetLogger().PushInfo(logger.GetFmtLog(anon_logger.CLogInfoUndecryptable, hash, proof, nil, conn))
			return
		}

		// check sender in f2f list
		if !node.GetListPubKeys().InPubKeys(sender) {
			node.GetLogger().PushWarn(logger.GetFmtLog(anon_logger.CLogWarnNotFriend, hash, proof, sender, conn))
			return
		}

		// get header of payload
		head := loadHead(pld.GetHead())

		// get session by payload head
		actionKey := newActionKey(sender, head.GetAction())
		action, ok := node.getAction(actionKey)
		if ok {
			node.GetLogger().PushInfo(logger.GetFmtLog(anon_logger.CLogInfoAction, hash, proof, sender, conn))
			action <- pld.GetBody()
			return
		}

		// get function by payload head
		f, ok := node.getRoute(head.GetRoute())
		if !ok || f == nil {
			node.GetLogger().PushWarn(logger.GetFmtLog(anon_logger.CLogWarnUnknownRoute, hash, proof, sender, conn))
			return
		}

		// response can be nil
		resp := f(node, sender, hash, pld.GetBody())
		if resp == nil {
			node.GetLogger().PushInfo(logger.GetFmtLog(anon_logger.CLogInfoWithoutResp, hash, proof, sender, conn))
			return
		}

		// create the message and put this to the queue
		if err := node.BroadcastPayload(sender, payload.NewPayload(pld.GetHead(), resp)); err != nil {
			node.GetLogger().PushErro(logger.GetFmtLog(anon_logger.CLogBaseEnqueueResp, hash, proof, sender, conn))
			return
		}

		node.GetLogger().PushInfo(logger.GetFmtLog(anon_logger.CLogBaseEnqueueResp, hash, proof, sender, conn))
	}
}

func (node *sNode) networkBroadcast(logger anon_logger.ILogger, msg message.IMessage) error {
	hash := msg.GetBody().GetHash()
	proof := msg.GetBody().GetProof()

	// redirect message to another nodes
	err := node.GetNetworkNode().BroadcastPayload(
		payload.NewPayload(
			node.GetSettings().GetNetworkMask(),
			msg.ToBytes(),
		),
	)
	if err != nil {
		node.fLogger.PushWarn(logger.GetFmtLog(anon_logger.CLogBaseBroadcast, hash, proof, nil, nil))
		// need continue (some of connections may be closed)
	}

	return nil
}

func (node *sNode) setRoute(head uint32, handle IHandlerF) {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	node.fHandleRoutes[head] = handle
}

func (node *sNode) getRoute(head uint32) (IHandlerF, bool) {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	f, ok := node.fHandleRoutes[head]
	return f, ok
}

func (node *sNode) getAction(actionKey string) (chan []byte, bool) {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	f, ok := node.fHandleActions[actionKey]
	return f, ok
}

func (node *sNode) setAction(actionKey string) {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	node.fHandleActions[actionKey] = make(chan []byte)
}

func (node *sNode) delAction(actionKey string) {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	delete(node.fHandleActions, actionKey)
}

func newActionKey(pubKey asymmetric.IPubKey, head uint32) string {
	pubKeyAddr := pubKey.Address().ToString()
	headString := fmt.Sprintf("%d", head)
	return fmt.Sprintf("%s-%s", pubKeyAddr, headString)
}

func mustBeUint32(v uint64) uint32 {
	if v > math.MaxUint32 {
		panic("v > math.MaxUint32")
	}
	return uint32(v)
}
