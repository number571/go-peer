package block

import (
	"github.com/number571/go-peer/cmd/union_blockchain/kernel"
	"github.com/number571/go-peer/cmd/union_blockchain/kernel/transaction"
)

type ISettings interface {
	GetCountTXs() uint64
	GetTransactionSettings() transaction.ISettings
}

type IBlock interface {
	Settings() ISettings

	PrevHash() []byte
	Transactions() []transaction.ITransaction

	kernel.IWrapper
	kernel.ISignifier
}
