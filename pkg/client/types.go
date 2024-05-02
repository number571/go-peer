package client

import (
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/payload"
)

type IClient interface {
	GetSettings() message.ISettings
	GetMessageLimit() uint64

	GetPubKey() asymmetric.IPubKey
	GetPrivKey() asymmetric.IPrivKey

	EncryptPayload(asymmetric.IPubKey, payload.IPayload64) (message.IMessage, error)
	DecryptMessage(message.IMessage) (asymmetric.IPubKey, payload.IPayload64, error)
}
