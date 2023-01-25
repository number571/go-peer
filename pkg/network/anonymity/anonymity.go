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
	"github.com/number571/go-peer/pkg/closer"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/friends"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/go-peer/pkg/types"

	"github.com/number571/go-peer/pkg/network/conn"
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
	fQueue         queue.IQueue
	fF2F           friends.IF2F
	fHandleRoutes  map[uint32]IHandlerF
	fHandleActions map[string]chan []byte
}

func NewNode(
	sett ISettings,
	log logger.ILogger,
	kvDB database.IKeyValueDB,
	nnode network.INode,
	queue queue.IQueue,
	f2f friends.IF2F,
) INode {
	return &sNode{
		fSettings:      sett,
		fLogger:        log,
		fKeyValueDB:    kvDB,
		fNetwork:       nnode,
		fQueue:         queue,
		fF2F:           f2f,
		fHandleRoutes:  make(map[uint32]IHandlerF),
		fHandleActions: make(map[string]chan []byte),
	}
}

func (node *sNode) Run() error {
	if err := node.runQueue(); err != nil {
		return err
	}

	node.Network().Handle(
		node.Settings().GetNetworkMask(),
		node.handleWrapper(),
	)

	return nil
}

func (node *sNode) Close() error {
	node.Network().Handle(node.Settings().GetNetworkMask(), nil)
	return closer.CloseAll([]types.ICloser{
		node.Queue(),
	})
}

func (node *sNode) Settings() ISettings {
	return node.fSettings
}

func (node *sNode) KeyValueDB() database.IKeyValueDB {
	return node.fKeyValueDB
}

func (node *sNode) Network() network.INode {
	return node.fNetwork
}

func (node *sNode) Queue() queue.IQueue {
	return node.fQueue
}

// Return f2f structure.
func (node *sNode) F2F() friends.IF2F {
	return node.fF2F
}

func (node *sNode) Handle(head uint32, handle IHandlerF) INode {
	node.setRoute(head, handle)
	return node
}

// Send message without response waiting.
func (node *sNode) Broadcast(recv asymmetric.IPubKey, pld payload.IPayload) error {
	if len(node.Network().Connections()) == 0 {
		return errors.New("length of connections = 0")
	}

	msg, err := node.Queue().Client().Encrypt(recv, pld)
	if err != nil {
		return err
	}

	return node.send(msg)
}

// Send message with response waiting.
// Payload head must be uint32.
func (node *sNode) Request(recv asymmetric.IPubKey, pld payload.IPayload) ([]byte, error) {
	if len(node.Network().Connections()) == 0 {
		return nil, errors.New("length of connections = 0")
	}

	headAction := uint32(random.NewStdPRNG().Uint64())
	headRoute := mustBeUint32(pld.Head())

	newPld := payload.NewPayload(
		joinHead(headAction, headRoute).Uint64(),
		pld.Body(),
	)

	actionKey := newActionKey(recv, headAction)

	node.setAction(actionKey)
	defer node.delAction(actionKey)

	if err := node.Broadcast(recv, newPld); err != nil {
		return nil, err
	}
	return node.recv(actionKey)
}

func (node *sNode) send(msg message.IMessage) error {
	for i := uint64(0); i <= node.Settings().GetRetryEnqueue(); i++ {
		if err := node.Queue().Enqueue(msg); err != nil {
			time.Sleep(node.Queue().Settings().GetDuration())
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
	case <-time.After(node.Settings().GetTimeWait()):
		return nil, errors.New("time is over")
	}
}

func (node *sNode) runQueue() error {
	if err := node.Queue().Run(); err != nil {
		return err
	}

	go func() {
		for {
			msg, ok := <-node.Queue().Dequeue()
			if !ok {
				break
			}

			var (
				hash   = msg.Body().Hash()
				proof  = msg.Body().Proof()
				pubKey = node.Queue().Client().PubKey()
			)

			if err := node.networkBroadcast(msg); err != nil {
				node.fLogger.Erro(fmtLog(cLogErroMiddleware, hash, proof, pubKey, nil))
				continue
			}

			node.fLogger.Info(fmtLog(cLogBaseBroadcast, hash, proof, pubKey, nil))
		}
	}()

	return nil
}

func (node *sNode) handleWrapper() network.IHandlerF {
	return func(_ network.INode, conn conn.IConn, reqBytes []byte) {
		msg := message.LoadMessage(
			reqBytes,
			node.Queue().Client().Settings().GetMessageSize(),
			node.Queue().Client().Settings().GetWorkSize(),
		)
		if msg == nil {
			node.fLogger.Warn(fmtLog(cLogWarnMessageNull, nil, 0, nil, conn))
			return
		}

		var (
			addr  = node.Queue().Client().PubKey().Address().String()
			hash  = msg.Body().Hash()
			proof = msg.Body().Proof()
		)

		hashDB := []byte(fmt.Sprintf("_hash_%X", hash))
		gotAddrs, err := node.KeyValueDB().Get(hashDB)

		// check already received data by hash
		hashIsExist := (err == nil)
		if hashIsExist && strings.Contains(string(gotAddrs), addr) {
			node.fLogger.Info(fmtLog(cLogInfoExist, hash, proof, nil, conn))
			return
		}

		// set hash to database
		updateAddrs := fmt.Sprintf("%s;%s", string(gotAddrs), addr)
		if err := node.KeyValueDB().Set(hashDB, []byte(updateAddrs)); err != nil {
			node.fLogger.Erro(fmtLog(cLogErroDatabaseSet, hash, proof, nil, conn))
			return
		}

		// do not send data if than already received
		if !hashIsExist {
			// broadcast message to network
			if err := node.networkBroadcast(msg); err != nil {
				node.fLogger.Erro(fmtLog(cLogErroMiddleware, hash, proof, nil, conn))
				return
			}
		}

		// try decrypt message
		sender, pld, err := node.Queue().Client().Decrypt(msg)
		if err != nil {
			node.fLogger.Info(fmtLog(cLogInfoUnencryptable, hash, proof, nil, conn))
			return
		}

		// check sender in f2f list
		if !node.F2F().InList(sender) {
			node.fLogger.Warn(fmtLog(cLogWarnNotFriend, hash, proof, sender, conn))
			return
		}

		// get header of payload
		head := loadHead(pld.Head())

		// get session by payload head
		actionKey := newActionKey(sender, head.GetAction())
		action, ok := node.getAction(actionKey)
		if ok {
			node.fLogger.Info(fmtLog(cLogInfoAction, hash, proof, sender, conn))
			action <- pld.Body()
			return
		}

		// get function by payload head
		f, ok := node.getRoute(head.GetRoute())
		if !ok || f == nil {
			node.fLogger.Warn(fmtLog(cLogWarnUnknownRoute, hash, proof, sender, conn))
			return
		}

		// response can be nil
		resp := f(node, sender, pld.Body())
		if resp == nil {
			node.fLogger.Info(fmtLog(cLogInfoWithoutResp, hash, proof, sender, conn))
			return
		}

		// create the message and put this to the queue
		if err := node.Broadcast(sender, payload.NewPayload(pld.Head(), resp)); err != nil {
			node.fLogger.Erro(fmtLog(cLogBaseEnqueueResp, hash, proof, sender, conn))
			return
		}

		node.fLogger.Info(fmtLog(cLogBaseEnqueueResp, hash, proof, sender, conn))
	}
}

func (node *sNode) networkBroadcast(msg message.IMessage) error {
	hash := msg.Body().Hash()
	proof := msg.Body().Proof()

	// redirect message to another nodes
	err := node.Network().Broadcast(payload.NewPayload(
		node.Settings().GetNetworkMask(),
		msg.Bytes(),
	))
	if err != nil {
		node.fLogger.Warn(fmtLog(cLogBaseBroadcast, hash, proof, nil, nil))
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
	pubKeyAddr := pubKey.Address().String()
	headString := fmt.Sprintf("%d", head)
	return fmt.Sprintf("%s-%s", pubKeyAddr, headString)
}

func mustBeUint32(v uint64) uint32 {
	if v > math.MaxUint32 {
		panic("v > math.MaxUint32")
	}
	return uint32(v)
}
