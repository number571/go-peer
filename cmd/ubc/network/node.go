package network

import (
	"sync"

	"github.com/number571/go-peer/cmd/ubc/kernel/chain"
	"github.com/number571/go-peer/modules/network"
	"github.com/number571/go-peer/modules/network/conn"
	"github.com/number571/go-peer/modules/payload"
)

var (
	_ INode = &sNode{}
)

type sNode struct {
	fMutex        sync.Mutex
	fNetwork      network.INode
	fPusher       IPusher
	fChain        chain.IChain
	fHandleRoutes map[uint64]IHandlerF
}

func NewNode(chain chain.IChain) INode {
	node := &sNode{
		fChain:        chain,
		fHandleRoutes: make(map[uint64]IHandlerF),
	}

	node.fPusher = NewPusher(node.fNetwork)
	node.fNetwork.Handle(cMaskNetw, node.handleWrapper())

	return node.
		Handle(cMaskPushBlock, handlePushBlock).
		Handle(cMaskPushTransaction, handlePushTransaction).
		Handle(cMaskLoadHeight, handleLoadHeight).
		Handle(cMaskLoadBlock, handleLoadBlock).
		Handle(cMaskLoadTransaction, handleLoadTransaction)
}

func (node *sNode) Network() network.INode {
	return node.fNetwork
}

func (node *sNode) Pusher() IPusher {
	return node.fPusher
}

func (node *sNode) Handle(head uint64, handle IHandlerF) INode {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	node.fHandleRoutes[head] = handle
	return node
}

func (node *sNode) handleWrapper() network.IHandlerF {
	return func(nnode network.INode, conn conn.IConn, npld payload.IPayload) {
		pld := payload.LoadPayload(npld.Body())

		// get function by payload head
		f, ok := node.getFunction(pld.Head())
		if !ok || f == nil {
			return
		}

		f(node, conn, pld)
	}
}

func (node *sNode) getFunction(head uint64) (IHandlerF, bool) {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	f, ok := node.fHandleRoutes[head]
	return f, ok
}
