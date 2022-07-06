package transaction

import "github.com/number571/go-peer/cmd/ubc/kernel"

type ITransaction interface {
	Payload() []byte

	kernel.IWrapper
	kernel.ISignifier
}
