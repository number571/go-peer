package hashing

import (
	"github.com/number571/go-peer/pkg/crypto"
	"github.com/number571/go-peer/pkg/types"
)

type IHasher interface {
	types.IConverter
	crypto.IParameter
}
