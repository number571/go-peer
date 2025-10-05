package client

import (
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

type IClient interface {
	IEncryptor
	IDecryptor

	GetPrivKey() asymmetric.IPrivKey

	GetMessageSize() uint64
	GetPayloadSize() uint64
}

type IDecryptor interface {
	DecryptMessage(asymmetric.IMapPubKeys, []byte) (asymmetric.IPubKey, []byte, error)
}

type IEncryptor interface {
	EncryptMessage(asymmetric.IPubKey, []byte) ([]byte, error)
}
