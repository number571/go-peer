package client

import (
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/quantum"
)

type IClient interface {
	GetSettings() message.ISettings
	GetMessageLimit() uint64

	GetPrivKeyChain() quantum.IPrivKeyChain

	EncryptMessage(quantum.IKEMPubKey, []byte) ([]byte, error)
	DecryptMessage([]byte) (quantum.ISignerPubKey, []byte, error)
}
