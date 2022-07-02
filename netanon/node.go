package netanon

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/crypto/hashing"
	"github.com/number571/go-peer/crypto/puzzle"
	"github.com/number571/go-peer/crypto/random"
	"github.com/number571/go-peer/local/client"
	"github.com/number571/go-peer/local/message"
	"github.com/number571/go-peer/local/payload"
	"github.com/number571/go-peer/local/routing"
	"github.com/number571/go-peer/network"
	"github.com/number571/go-peer/settings"
)

var (
	_ INode = &sNode{}
)

type sNode struct {
	fMutex         sync.Mutex
	fPseudoClient  asymmetric.IPrivKey
	fClient        client.IClient
	fNetwork       network.INode
	fRouterF       IRouterF
	fHandleRoutes  map[uint32]IHandlerF
	fHandleActions map[uint32]chan []byte
	fF2F           iF2F
	fOnline        iOnline
	// fChecker       iChecker
}

func NewNode(client client.IClient) INode {
	sett := client.Settings()
	node := &sNode{
		fPseudoClient:  asymmetric.NewRSAPrivKey(client.PrivKey().Size()),
		fClient:        client,
		fNetwork:       network.NewNode(sett),
		fRouterF:       func() []asymmetric.IPubKey { return nil },
		fHandleRoutes:  make(map[uint32]IHandlerF),
		fHandleActions: make(map[uint32]chan []byte),
		fF2F:           newF2F(),
	}

	// recurrent structures
	{
		node.fOnline = newOnline(node)
		// node.fChecker = newChecker(node)
	}

	node.fNetwork.Handle(
		sett.Get(settings.CMaskNetw),
		node.handleWrapper(),
	)

	return node
}

func (node *sNode) Close() error {
	return node.Network().Close()
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

// Return f2f structure.
func (node *sNode) F2F() iF2F {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	return node.fF2F
}

// Return online structure.
func (node *sNode) Online() iOnline {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	// return node.fOnline
	return nil
}

func (node *sNode) WithRouter(router IRouterF) INode {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	node.fRouterF = router
	return node
}

func (node *sNode) Handle(head uint64, handle IHandlerF) INode {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	// used only 32bit from 64bit number
	node.fHandleRoutes[loadHead(head).Routes()] = handle
	return node
}

func (node *sNode) Broadcast(msg message.IMessage) error {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	return node.fNetwork.Broadcast(payload.NewPayload(
		node.fClient.Settings().Get(settings.CMaskNetw),
		msg.Bytes(),
	))
}

// Send message by public key of receiver.
func (node *sNode) Request(recv asymmetric.IPubKey, pl payload.IPayload) ([]byte, error) {
	return node.doRequest(
		recv,
		pl,
		node.fRouterF,
		node.Client().Settings().Get(settings.CSizeRtry),
		node.Client().Settings().Get(settings.CTimeWait),
	)
}

func (node *sNode) doRequest(recv asymmetric.IPubKey, pl payload.IPayload, fRoute IRouterF, retryNum, timeWait uint64) ([]byte, error) {
	if len(node.Network().Connections()) == 0 {
		return nil, errors.New("length of connections = 0")
	}

	headRoutes := loadHead(pl.Head()).Routes()
	headAction := loadHead(random.NewStdPRNG().Uint64()).Actions()

	pl = payload.NewPayload(
		joinHead(headRoutes, headAction).Uint64(),
		pl.Body(),
	)

	route := routing.NewRoute(recv).WithRedirects(node.fPseudoClient, fRoute())
	routeMsg := node.Client().Encrypt(route, pl)
	if routeMsg == nil {
		return nil, errors.New("psender is nil with not nil routes")
	}

	node.setAction(headAction)
	defer node.delAction(headAction)

	for counter := uint64(0); counter <= retryNum; counter++ {
		node.Broadcast(routeMsg)
		resp, err := node.recv(headAction, timeWait)
		if err != nil {
			return nil, err
		}
		if resp == nil {
			continue
		}
		return resp, nil
	}

	return nil, errors.New("time is over")
}

func (node *sNode) recv(head uint32, timeOut uint64) ([]byte, error) {
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
	case <-time.After(time.Duration(timeOut) * time.Second):
		return nil, nil
	}
}

func (node *sNode) handleWrapper() network.IHandlerF {
	return func(nnode network.INode, _ network.IConn, npld payload.IPayload) {
		msg := node.initialCheck(message.LoadMessage(npld.Body()))
		if msg == nil {
			return
		}

		// TODO: time.sleep(random)
		// TODO: f2f mode
		// TODO: random send pseudo message
		// TODO: request/response message

		for {
			node.Broadcast(msg)

			// try decrypt message
			sender, pld := node.Client().Decrypt(msg)
			if pld == nil {
				return
			}

			switch pld.Head() {
			case node.Client().Settings().Get(settings.CMaskRout):
				// redirect decrypt message
				msg = message.LoadMessage(pld.Body())
				if msg != nil {
					break // switch
				}
				return
			default:
				// check f2f mode and sender in f2f list
				if node.fF2F.Status() && !node.fF2F.InList(sender) {
					return
				}

				// get session by payload head
				action, ok := node.getAction(
					loadHead(pld.Head()).Actions(),
				)
				if ok {
					action <- pld.Body()
					return
				}

				// get function by payload head
				f, ok := node.getFunction(
					loadHead(pld.Head()).Routes(),
				)
				if !ok || f == nil {
					return
				}

				fmt.Println(string(pld.Body()))

				// send response
				// problem - withredirects?
				node.Broadcast(node.Client().Encrypt(
					routing.NewRoute(sender).WithRedirects(
						node.fPseudoClient,
						node.fRouterF(),
					),
					payload.NewPayload(pld.Head(), f(node, sender, pld)),
				))
				return
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

	diff := node.fClient.Settings().Get(settings.CSizeWork)
	puzzle := puzzle.NewPoWPuzzle(diff)
	if !puzzle.Verify(msg.Body().Hash(), msg.Body().Proof()) {
		return nil
	}

	return msg
}

func (node *sNode) getFunction(head uint32) (IHandlerF, bool) {
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
