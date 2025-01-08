package asymmetric

import (
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/types"
)

type IMapPubKeys interface {
	GetPubKey(string) IPubKey
	DelPubKey(IPubKey)
	SetPubKey(IPubKey)
}

type IPrivKey interface {
	types.IConverter

	GetPubKey() IPubKey
	GetKEMPrivKey() IKEMPrivKey
	GetDSAPrivKey() IDSAPrivKey
}

type IPubKey interface {
	types.IConverter
	GetHasher() hashing.IHasher

	GetKEMPubKey() IKEMPubKey
	GetDSAPubKey() IDSAPubKey
}

type IKEMPrivKey interface {
	ToBytes() []byte
	GetPubKey() IKEMPubKey

	Decapsulate([]byte) ([]byte, error)
}

type IKEMPubKey interface {
	ToBytes() []byte

	Encapsulate() ([]byte, []byte, error)
}

type IDSAPrivKey interface {
	ToBytes() []byte
	GetPubKey() IDSAPubKey

	SignBytes([]byte) []byte
}

type IDSAPubKey interface {
	ToBytes() []byte

	VerifyBytes([]byte, []byte) bool
}
