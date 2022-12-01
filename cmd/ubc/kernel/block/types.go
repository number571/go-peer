package block

import (
	"github.com/number571/go-peer/cmd/ubc/kernel"
	"github.com/number571/go-peer/cmd/ubc/kernel/transaction"
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
