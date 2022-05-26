package hmc

import (
	"github.com/number571/go-peer/crypto"
)

type IClient interface {
	Size() (uint64, error)
	Load(uint64) ([]byte, error)
	Push(crypto.IPubKey, []byte) error
}
