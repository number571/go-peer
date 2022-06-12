package symmetric

import "github.com/number571/go-peer/crypto"

type ICipher interface {
	crypto.IEncrypter
	crypto.IDecrypter
	crypto.IConverter
}
