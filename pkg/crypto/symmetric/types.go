package symmetric

import (
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/types"
)

type ICipher interface {
	types.IConverter
	GetHasher() hashing.IHasher

	EncryptBytes(pMsg []byte) []byte
	DecryptBytes(pMsg []byte) []byte
}
