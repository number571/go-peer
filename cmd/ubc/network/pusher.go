package network

import (
	"sync"

	"github.com/number571/go-peer/cmd/ubc/kernel/block"
	"github.com/number571/go-peer/cmd/ubc/kernel/transaction"
	"github.com/number571/go-peer/local/payload"
	"github.com/number571/go-peer/network"
)

var (
	_ IPusher = &sPusher{}
)

type sPusher struct {
	fMutex   sync.Mutex
	fNetwork network.INode
}

func NewPusher(node network.INode) IPusher {
	return &sPusher{
		fNetwork: node,
	}
}

func (pusher *sPusher) Block(block block.IBlock) error {
	pusher.fMutex.Lock()
	defer pusher.fMutex.Unlock()

	pld := payload.NewPayload(
		cMaskPushBlock,
		block.Bytes(),
	)
	return pusher.fNetwork.Broadcast(payload.NewPayload(
		cMaskNetw,
		pld.Bytes(),
	))
}

func (pusher *sPusher) Transaction(transaction transaction.ITransaction) error {
	pusher.fMutex.Lock()
	defer pusher.fMutex.Unlock()

	pld := payload.NewPayload(
		cMaskPushTransaction,
		transaction.Bytes(),
	)
	return pusher.fNetwork.Broadcast(payload.NewPayload(
		cMaskNetw,
		pld.Bytes(),
	))
}
