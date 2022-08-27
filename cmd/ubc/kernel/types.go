package kernel

import (
	"github.com/number571/go-peer/modules/crypto/asymmetric"
)

type IWrapper interface {
	Bytes() []byte
	String() string
}

type IHasher interface {
	Hash() []byte
	IsValid() bool
}

type ISignifier interface {
	Sign() []byte
	Validator() asymmetric.IPubKey

	IHasher
}
