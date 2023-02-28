package client

import (
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/payload"
)

type IClient interface {
	GetSettings() ISettings

	GetPubKey() asymmetric.IPubKey
	GetPrivKey() asymmetric.IPrivKey

	EncryptPayload(asymmetric.IPubKey, payload.IPayload) (message.IMessage, error)
	DecryptMessage(message.IMessage) (asymmetric.IPubKey, payload.IPayload, error)
}

type ISettings interface {
	GetMessageSize() uint64
	GetWorkSize() uint64
}
