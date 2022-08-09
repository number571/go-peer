package netanon

import (
	"errors"
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
	"github.com/number571/go-peer/routing"
	"github.com/number571/go-peer/settings"
	"github.com/number571/go-peer/utils"
)

var (
	_ INode = &sNode{}
)

type sNode struct {
	fMutex         sync.Mutex
	fSettings      ISettings
	fClient        client.IClient
	fPseudo        asymmetric.IPrivKey
	fNetwork       network.INode
	fQueue         queue.IQueue
	fF2F           friends.IF2F
	fRouterF       IRouterF
	fHandleRoutes  map[uint32]IHandlerF
	fHandleActions map[uint32]chan []byte
}

func NewNode(
	sett ISettings,
	client client.IClient,
	nnode network.INode,
	queue queue.IQueue,
	f2f friends.IF2F,
	route IRouterF,
) INode {
	// queue.NewQueue(client, settings.CSizeDefaultCap, sett.Get(settings.CSizePerd))
	newKey := asymmetric.NewRSAPrivKey(client.PrivKey().Size())
	node := &sNode{
		fSettings:      sett,
		fClient:        client,
		fPseudo:        newKey,
		fNetwork:       nnode,
		fQueue:         queue,
		fF2F:           f2f,
		fRouterF:       route,
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

func (node *sNode) WithRouter(router IRouterF) INode {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	node.fRouterF = router
	return node
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
func (node *sNode) Request(recv asymmetric.IPubKey, pl payload.IPayload) ([]byte, error) {
	return node.doRequest(
		recv,
		pl,
		node.Settings().GetTimeWait(),
	)
}

func (node *sNode) doRequest(recv asymmetric.IPubKey, pl payload.IPayload, timeWait time.Duration) ([]byte, error) {
	if len(node.Network().Connections()) == 0 {
		return nil, errors.New("length of connections = 0")
	}

	headRoutes := utils.MustBeUint32(pl.Head())
	headAction := uint32(random.NewStdPRNG().Uint64())

	pl = payload.NewPayload(
		joinHead(headRoutes, headAction).Uint64(),
		pl.Body(),
	)

	route := routing.NewRoute(recv).WithRedirects(node.fPseudo, node.fRouterF())
	routeMsg := node.Client().Encrypt(route, pl)
	if routeMsg == nil {
		return nil, errors.New("psender is nil with not nil routes")
	}

	node.setAction(headAction)
	defer node.delAction(headAction)

	for i := uint64(0); i < node.Settings().GetRetryEnqueue(); i++ {
		if err := node.Queue().Enqueue(settings.CNull, routeMsg); err != nil {
			time.Sleep(node.Queue().Settings().GetDuration())
			continue
		}
		break
	}
	return node.recv(headAction, timeWait)
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

		// msg can be decrypted for next route
		// for this need use Broadcast from netanon
		node.Broadcast(msg)

		// try decrypt message
		sender, pld := node.Client().Decrypt(msg)
		if pld == nil {
			return
		}

		head := pld.Head()
		switch head {
		case settings.CMaskRoute:
			// redirect decrypt message
			msg = message.LoadMessage(pld.Body())
			if msg == nil {
				return
			}
			for i := uint64(0); i < node.Settings().GetRetryEnqueue(); i++ {
				if err := node.Queue().Enqueue(settings.CNull, msg); err != nil {
					time.Sleep(node.Queue().Settings().GetDuration())
					continue
				}
				break
			}
		default:
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

			msg := node.Client().Encrypt(
				routing.NewRoute(sender).WithRedirects(
					node.fPseudo,
					node.fRouterF(),
				),
				payload.NewPayload(head, resp),
			)

			// send response with two append enqueue
			for i := uint64(0); i < node.Settings().GetRetryEnqueue(); i++ {
				err := node.Queue().Enqueue(node.Settings().GetResponsePeriod(), msg)
				if err != nil {
					time.Sleep(node.Queue().Settings().GetDuration())
					continue
				}
				break
			}
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
