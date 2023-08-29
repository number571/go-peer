package client

import (
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/payload"
)

type IClient interface {
	GetSettings() message.ISettings

	GetPubKey() asymmetric.IPubKey
	GetPrivKey() asymmetric.IPrivKey

	GetMessageLimit() uint64

	EncryptPayload(asymmetric.IPubKey, payload.IPayload) (message.IMessage, error)
	DecryptMessage(message.IMessage) (asymmetric.IPubKey, payload.IPayload, error)
}
