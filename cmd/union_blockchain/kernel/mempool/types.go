package mempool

import (
	"github.com/number571/go-peer/cmd/union_blockchain/kernel/block"
	"github.com/number571/go-peer/cmd/union_blockchain/kernel/transaction"
	"github.com/number571/go-peer/pkg/types"
)

type ISettings interface {
	GetCountTXs() uint64
	GetBlockSettings() block.ISettings
}

type IMempool interface {
	Settings() ISettings

	Height() uint64
	Transaction([]byte) transaction.ITransaction

	Push(transaction.ITransaction)
	Pop() []transaction.ITransaction

	Delete([]byte)
	Clear()
	types.ICloser
}
