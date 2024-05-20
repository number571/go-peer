package client

import (
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

type IClient interface {
	GetSettings() message.ISettings
	GetMessageLimit() uint64

	GetPubKey() asymmetric.IPubKey
	GetPrivKey() asymmetric.IPrivKey

	EncryptMessage(asymmetric.IPubKey, []byte) ([]byte, error)
	DecryptMessage([]byte) (asymmetric.IPubKey, []byte, error)
}
