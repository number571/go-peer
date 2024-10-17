package asymmetric

import (
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/types"
)

type IListPubKeyChains interface {
	AllPubKeyChains() []IPubKeyChain
	GetPubKeyChain(ISignPubKey) (IPubKeyChain, bool)
	AddPubKeyChain(IPubKeyChain)
	DelPubKeyChain(IPubKeyChain)
}

type IPrivKeyChain interface {
	types.IConverter

	GetPubKeyChain() IPubKeyChain
	GetKEncPrivKey() IKEncPrivKey
	GetSignPrivKey() ISignPrivKey
}

type IPubKeyChain interface {
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
