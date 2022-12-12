package symmetric

import "github.com/number571/go-peer/pkg/crypto"

type ICipher interface {
	crypto.IEncrypter
	crypto.IDecrypter
	crypto.IConverter
}
