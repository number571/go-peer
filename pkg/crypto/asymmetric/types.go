package asymmetric

import "github.com/number571/go-peer/pkg/crypto"

type iAddress interface {
	crypto.IConverter
}

type IPubKey interface {
	crypto.IEncrypter
	crypto.IConverter

	Address() iAddress
	Verify([]byte, []byte) bool
}

type IPrivKey interface {
	crypto.IDecrypter
	crypto.IConverter

	Sign([]byte) []byte
	PubKey() IPubKey
}
