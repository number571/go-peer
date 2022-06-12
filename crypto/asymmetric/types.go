package asymmetric

import "github.com/number571/go-peer/crypto"

type IPubKey interface {
	crypto.IEncrypter
	crypto.IConverter
	Address() string
	Verify([]byte, []byte) bool
}

type IPrivKey interface {
	crypto.IDecrypter
	crypto.IConverter
	Sign([]byte) []byte
	PubKey() IPubKey
}
