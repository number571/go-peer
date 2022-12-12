package asymmetric

import "github.com/number571/go-peer/pkg/crypto"

type iFormatter interface {
	Format() string
}

type iAddress interface {
	crypto.IConverter
}

type IPubKey interface {
	iFormatter

	crypto.IEncrypter
	crypto.IConverter

	Address() iAddress
	Verify([]byte, []byte) bool
}

type IPrivKey interface {
	iFormatter

	crypto.IDecrypter
	crypto.IConverter

	Sign([]byte) []byte
	PubKey() IPubKey
}
