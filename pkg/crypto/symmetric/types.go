package symmetric

import "github.com/number571/go-peer/pkg/types"

type IListCiphers interface {
	Get() []ICipher
	Add(ICipher) bool
	Del(ICipher) bool
}

type ICipher interface {
	types.IConverter

	EncryptBytes(pMsg []byte) []byte
	DecryptBytes(pMsg []byte) []byte
}
