package hashing

import "github.com/number571/go-peer/modules/crypto"

type IHasher interface {
	crypto.IConverter
}
