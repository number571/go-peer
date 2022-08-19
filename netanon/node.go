package netanon

import (
	"errors"
	"math"
	"sync"
	"time"

	"github.com/number571/go-peer/client"
	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/crypto/hashing"
	"github.com/number571/go-peer/crypto/puzzle"
	"github.com/number571/go-peer/crypto/random"
	"github.com/number571/go-peer/friends"
	"github.com/number571/go-peer/message"
	"github.com/number571/go-peer/network"
	"github.com/number571/go-peer/payload"
	"github.com/number571/go-peer/queue"
	"github.com/number571/go-peer/settings"

	payload_adapter "github.com/number571/go-peer/netanon/adapters/payload"
)

var (
	_ INode = &sNode{}
)

type sNode struct {
	fMutex         sync.Mutex
	fSettings      ISettings
	fClient        client.IClient
	fNetwork       network.INode
	fQueue         queue.IQueue
	fF2F           friends.IF2F
	fHandleRoutes  map[uint32]IHandlerF
	fHandleActions map[uint32]chan []byte
}

func NewNode(
	sett ISettings,
	client client.IClient,
	nnode network.INode,
	queue queue.IQueue,
	f2f friends.IF2F,
) INode {
	node := &sNode{
		fSettings:      sett,
		fClient:        client,
		fNetwork:       nnode,
		fQueue:         queue,
		fF2F:           f2f,
		fHandleRoutes:  make(map[uint32]IHandlerF),
		fHandleActions: make(map[uint32]chan []byte),
	}
	node.Queue().Start()
	node.Network().Handle(settings.CMaskNetwork, node.handleWrapper())
	return node
}

func (node *sNode) Close() error {
	var lerr error
	if err := node.Network().Close(); err != nil {
		lerr = err
	}
	if err := node.Queue().Close(); err != nil {
		lerr = err
	}
	return lerr
}

func (node *sNode) Settings() ISettings {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	return node.fSettings
}

func (node *sNode) Client() client.IClient {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	return node.fClient
}

func (node *sNode) Network() network.INode {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	return node.fNetwork
}

func (node *sNode) Queue() queue.IQueue {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	return node.fQueue
}

// Return f2f structure.
func (node *sNode) F2F() friends.IF2F {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	return node.fF2F
}

func (node *sNode) Handle(head uint32, handle IHandlerF) INode {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	node.fHandleRoutes[head] = handle
	return node
}

func (node *sNode) Broadcast(msg message.IMessage) error {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	return node.fNetwork.Broadcast(payload.NewPayload(
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

	msg := node.Client().Encrypt(recv, pl)

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
			msg := <-node.Queue().Dequeue()
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

		// try decrypt message
		sender, pld := node.Client().Decrypt(msg)
		if pld == nil {
			return
		}

		head := pld.Head()

		// check f2f mode and sender in f2f list
		if node.F2F().Status() && !node.F2F().InList(sender) {
			return
		}

		// get session by payload head
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

		respMsg := node.Client().Encrypt(
			sender,
			payload.NewPayload(head, resp),
		)

		// send response with two append enqueue
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
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	if msg == nil {
		return nil
	}

	if len(msg.Body().Hash()) != hashing.GSHA256Size {
		return nil
	}

	diff := node.fClient.Settings().GetWorkSize()
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
