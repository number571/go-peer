package asymmetric

import (
	"github.com/number571/go-peer/pkg/crypto"
	"github.com/number571/go-peer/pkg/types"
)

type IAddress interface {
	types.IConverter
	types.IParameter
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
	types.IParameter

	GetAddress() IAddress
	VerifyBytes([]byte, []byte) bool
}

type IPrivKey interface {
	crypto.IDecrypter
	types.IConverter
	types.IParameter

	SignBytes([]byte) []byte
	GetPubKey() IPubKey
}

type IEphPubKey interface {
	types.IConverter
	types.IParameter
}

type IEphPrivKey interface {
	types.IConverter
	types.IParameter

	GetSharedKey(IEphPubKey) ([]byte, error)
	GetPubKey() IEphPubKey
}
