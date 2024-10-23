package client

import (
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

type IClient interface {
	GetMessageSize() uint64
	GetPayloadLimit() uint64

	GetPrivKey() asymmetric.IPrivKey

	EncryptMessage(asymmetric.IPubKey, []byte) ([]byte, error)
	DecryptMessage(asymmetric.IMapPubKeys, []byte) (asymmetric.IPubKey, []byte, error)
}
