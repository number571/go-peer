package asymmetric

import (
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/types"
)

type IListPubKeys interface {
	AllPubKeys() []IPubKey
	GetPubKey(ISignPubKey) (IPubKey, bool)
	AddPubKey(IPubKey)
	DelPubKey(IPubKey)
}

type IPrivKey interface {
	types.IConverter

	GetPubKey() IPubKey
	GetKEncPrivKey() IKEncPrivKey
	GetSignPrivKey() ISignPrivKey
}

type IPubKey interface {
	types.IConverter
	GetHasher() hashing.IHasher

	GetKEncPubKey() IKEncPubKey
	GetSignPubKey() ISignPubKey
}

type IKEncPrivKey interface {
	ToBytes() []byte
	GetPubKey() IKEncPubKey

	Decapsulate([]byte) ([]byte, error)
}

type IKEncPubKey interface {
	ToBytes() []byte

	Encapsulate() ([]byte, []byte, error)
}

type ISignPrivKey interface {
	ToBytes() []byte
	GetPubKey() ISignPubKey

	SignBytes([]byte) []byte
}

type ISignPubKey interface {
	ToBytes() []byte

	VerifyBytes([]byte, []byte) bool
}
