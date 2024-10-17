package client

import (
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

type IClient interface {
	GetSettings() message.ISettings
	GetMessageLimit() uint64

	GetPrivKeyChain() asymmetric.IPrivKeyChain

	EncryptMessage(asymmetric.IKEncPubKey, []byte) ([]byte, error)
	DecryptMessage([]byte) (asymmetric.ISignPubKey, []byte, error)
}
