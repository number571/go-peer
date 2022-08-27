package asymmetric

import "github.com/number571/go-peer/modules/crypto"

type IPubKey interface {
	crypto.IEncrypter
	crypto.IConverter
	Address() iAddress
	Verify([]byte, []byte) bool
}

type iAddress interface {
	crypto.IConverter
}

type IPrivKey interface {
	crypto.IDecrypter
	crypto.IConverter
	Sign([]byte) []byte
	PubKey() IPubKey
}
