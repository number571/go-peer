package asymmetric

import (
	"github.com/number571/go-peer/pkg/crypto"
	"github.com/number571/go-peer/pkg/types"
)

type IAddress interface {
	types.IConverter
	crypto.IParameter
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
	crypto.IParameter

	GetAddress() IAddress
	VerifyBytes([]byte, []byte) bool
}

type IPrivKey interface {
	crypto.IDecrypter
	types.IConverter
	crypto.IParameter

	SignBytes([]byte) []byte
	GetPubKey() IPubKey
}
