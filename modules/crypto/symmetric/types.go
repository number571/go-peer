package symmetric

import "github.com/number571/go-peer/modules/crypto"

type ICipher interface {
	crypto.IEncrypter
	crypto.IDecrypter
	crypto.IConverter
}
