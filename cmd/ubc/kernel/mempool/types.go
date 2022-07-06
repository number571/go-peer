package mempool

import "github.com/number571/go-peer/cmd/ubc/kernel/transaction"

type IMempool interface {
	Height() uint64
	Transaction([]byte) transaction.ITransaction

	Push(transaction.ITransaction)
	Pop() []transaction.ITransaction

	Delete([]byte)
	Clear()
	Close() error
}
