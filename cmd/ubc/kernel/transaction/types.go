package transaction

import "github.com/number571/go-peer/cmd/ubc/kernel"

type ISettings interface {
	GetMaxSize() uint64
}

type ITransaction interface {
	Settings() ISettings
	Payload() []byte

	kernel.IWrapper
	kernel.ISignifier
}
