package hashing

import "github.com/number571/go-peer/crypto"

type IHasher interface {
	crypto.IConverter
}
