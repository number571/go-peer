package network

import (
	"github.com/number571/go-peer/cmd/ubc/kernel/block"
	"github.com/number571/go-peer/cmd/ubc/kernel/transaction"
	"github.com/number571/go-peer/local/payload"
	"github.com/number571/go-peer/network"
)

type IHandlerF func(INode, network.IConn, payload.IPayload)

type INode interface {
	Pusher() IPusher

	Network() network.INode
	Handle(uint64, IHandlerF) INode
}

type IPusher interface {
	Block(block.IBlock) error
	Transaction(transaction.ITransaction) error
}

type ILoader interface {
	Height() (uint64, error)
	Block(uint64) (block.IBlock, error)
	Transaction([]byte) (transaction.ITransaction, error)
}
