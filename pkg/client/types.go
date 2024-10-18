package client

import (
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

type IClient interface {
	GetMessageSize() uint64
	GetPayloadLimit() uint64

	GetPrivKey() asymmetric.IPrivKey

	EncryptMessage(asymmetric.IKEMPubKey, []byte) ([]byte, error)
	DecryptMessage([]byte) (asymmetric.IDSAPubKey, []byte, error)
}
