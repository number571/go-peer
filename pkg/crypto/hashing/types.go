package hashing

import "github.com/number571/go-peer/pkg/crypto"

type IHasher interface {
	crypto.IConverter
}
