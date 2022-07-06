package block

import (
	"github.com/number571/go-peer/cmd/ubc/kernel"
	"github.com/number571/go-peer/cmd/ubc/kernel/transaction"
)

type IBlock interface {
	PrevHash() []byte
	Transactions() []transaction.ITransaction

	kernel.IWrapper
	kernel.ISignifier
}
