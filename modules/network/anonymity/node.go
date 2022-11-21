package anonymity

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/number571/go-peer/modules"
	"github.com/number571/go-peer/modules/closer"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/crypto/hashing"
	"github.com/number571/go-peer/modules/crypto/puzzle"
	"github.com/number571/go-peer/modules/crypto/random"
	"github.com/number571/go-peer/modules/friends"
	"github.com/number571/go-peer/modules/message"
	"github.com/number571/go-peer/modules/network"
	"github.com/number571/go-peer/modules/payload"
	"github.com/number571/go-peer/modules/queue"
	"github.com/number571/go-peer/modules/storage/database"
	"github.com/number571/go-peer/settings"

	payload_adapter "github.com/number571/go-peer/modules/network/anonymity/adapters/payload"
)

var (
	_ INode = &sNode{}
)

type sNode struct {
	fMutex         sync.Mutex
	fSettings      ISettings
	fKeyValueDB    database.IKeyValueDB
	fNetwork       network.INode
	fQueue         queue.IQueue
	fF2F           friends.IF2F
	fHandleRoutes  map[uint32]IHandlerF
	fHandleActions map[uint32]chan []byte
}

func NewNode(
	sett ISettings,
	kvDB database.IKeyValueDB,
	nnode network.INode,
	queue queue.IQueue,
	f2f friends.IF2F,
) INode {
	return &sNode{
		fSettings:      sett,
		fKeyValueDB:    kvDB,
		fNetwork:       nnode,
		fQueue:         queue,
		fF2F:           f2f,
		fHandleRoutes:  make(map[uint32]IHandlerF),
		fHandleActions: make(map[uint32]chan []byte),
	}
}

func (node *sNode) Run() error {
	if err := node.Queue().Run(); err != nil {
		return err
	}
	node.Network().Handle(settings.CMaskNetwork, node.handleWrapper())
	return nil
}

func (node *sNode) Close() error {
	node.Network().Handle(settings.CMaskNetwork, nil)
	return closer.CloseAll([]modules.ICloser{
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
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	node.fHandleRoutes[head] = handle
	return node
}

func (node *sNode) Broadcast(msg message.IMessage) error {
	return node.Network().Broadcast(payload.NewPayload(
		settings.CMaskNetwork,
		msg.Bytes(),
	))
}

// Send message by public key of receiver.
func (node *sNode) Request(recv asymmetric.IPubKey, pl payload_adapter.IPayload) ([]byte, error) {
	if len(node.Network().Connections()) == 0 {
		return nil, errors.New("length of connections = 0")
	}

	headRoutes := mustBeUint32(pl.Head())
	headAction := uint32(random.NewStdPRNG().Uint64())

	pl = payload.NewPayload(
		joinHead(headRoutes, headAction).Uint64(),
		pl.Body(),
	)

	msg, err := node.Queue().Client().Encrypt(recv, pl)
	if err != nil {
		return nil, err
	}

	node.setAction(headAction)
	defer node.delAction(headAction)

	for i := uint64(0); i <= node.Settings().GetRetryEnqueue(); i++ {
		if err := node.Queue().Enqueue(msg); err != nil {
			time.Sleep(node.Queue().Settings().GetDuration())
			continue
		}
		break
	}
	return node.recv(headAction, node.Settings().GetTimeWait())
}

func (node *sNode) recv(head uint32, timeOut time.Duration) ([]byte, error) {
	action, ok := node.getAction(head)
	if !ok {
		return nil, errors.New("action undefined")
	}
	select {
	case result, opened := <-action:
		if !opened {
			return nil, errors.New("chan is closed")
		}
		return result, nil
	case <-time.After(timeOut):
		return nil, errors.New("time is over")
	}
}

func (node *sNode) handleWrapper() network.IHandlerF {
	go func() {
		for {
			msg, ok := <-node.Queue().Dequeue()
			if !ok {
				break
			}
			node.Broadcast(msg)
		}
	}()

	return func(nnode network.INode, _ network.IConn, npld payload.IPayload) {
		msg := node.initialCheck(message.LoadMessage(npld.Body()))
		if msg == nil {
			return
		}

		// redirect to another nodes
		nnode.Broadcast(npld)
		client := node.Queue().Client()

		// try decrypt message
		sender, pld, err := client.Decrypt(msg)
		if err != nil {
			return
		}

		// check sender in f2f list
		if !node.F2F().InList(sender) {
			return
		}

		// check already received data by hash
		hash := []byte(fmt.Sprintf("recv_hash_%X", msg.Body().Hash()))
		if _, err := node.KeyValueDB().Get(hash); err == nil {
			return
		}
		node.KeyValueDB().Set(hash, []byte{})

		// get session by payload head
		head := pld.Head()
		action, ok := node.getAction(
			loadHead(head).Actions(),
		)
		if ok {
			action <- pld.Body()
			return
		}

		// get function by payload head
		f, ok := node.getRoute(
			loadHead(head).Routes(),
		)
		if !ok || f == nil {
			return
		}

		resp := f(node, sender, pld)
		if resp == nil {
			return
		}

		respMsg, err := client.Encrypt(
			sender,
			payload.NewPayload(head, resp),
		)
		if err != nil {
			panic(err)
		}

		for i := uint64(0); i <= node.Settings().GetRetryEnqueue(); i++ {
			err := node.Queue().Enqueue(respMsg)
			if err != nil {
				time.Sleep(node.Queue().Settings().GetDuration())
				continue
			}
			break
		}
	}
}

func (node *sNode) initialCheck(msg message.IMessage) message.IMessage {
	if msg == nil {
		return nil
	}

	if len(msg.Body().Hash()) != hashing.CSHA256Size {
		return nil
	}

	diff := node.Queue().Client().Settings().GetWorkSize()
	puzzle := puzzle.NewPoWPuzzle(diff)
	if !puzzle.Verify(msg.Body().Hash(), msg.Body().Proof()) {
		return nil
	}

	return msg
}

func (node *sNode) getRoute(head uint32) (IHandlerF, bool) {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	f, ok := node.fHandleRoutes[head]
	return f, ok
}

func (node *sNode) getAction(head uint32) (chan []byte, bool) {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	f, ok := node.fHandleActions[head]
	return f, ok
}

func (node *sNode) setAction(head uint32) {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	node.fHandleActions[head] = make(chan []byte)
}

func (node *sNode) delAction(head uint32) {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	delete(node.fHandleActions, head)
}

func mustBeUint32(v uint64) uint32 {
	if v > math.MaxUint32 {
		panic("v > math.MaxUint32")
	}
	return uint32(v)
}
