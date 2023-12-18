package asymmetric

import (
	"github.com/number571/go-peer/pkg/crypto"
	"github.com/number571/go-peer/pkg/types"
)

type IAddress interface {
	types.IConverter
	GetSize() uint64
}

type IListPubKeys interface {
	InPubKeys(IPubKey) bool
	GetPubKeys() []IPubKey
	AddPubKey(IPubKey)
	DelPubKey(IPubKey)
}

type IPubKey interface {
	crypto.IEncrypter
	types.IConverter
	GetSize() uint64

	GetAddress() IAddress
	VerifyBytes([]byte, []byte) bool
}

type IPrivKey interface {
	crypto.IDecrypter
	types.IConverter
	GetSize() uint64

	SignBytes([]byte) []byte
	GetPubKey() IPubKey
}
