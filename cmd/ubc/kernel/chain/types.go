package chain

import (
	"github.com/number571/go-peer/cmd/ubc/kernel/block"
	"github.com/number571/go-peer/cmd/ubc/kernel/mempool"
	"github.com/number571/go-peer/cmd/ubc/kernel/transaction"
)

type IChain interface {
	Accept(block.IBlock) bool
	Merge([]transaction.ITransaction) bool
	Rollback(uint64) bool

	Height() uint64
	Transaction([]byte) transaction.ITransaction
	Block(uint64) block.IBlock

	Mempool() mempool.IMempool
	Close() error
}
