package client

import (
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/payload"
)

type IClient interface {
	Settings() ISettings

	iKeeper
	iCipher
}

type ISettings interface {
	GetMessageSize() uint64
	GetWorkSize() uint64
}

type iKeeper interface {
	PubKey() asymmetric.IPubKey
	PrivKey() asymmetric.IPrivKey
}

type iCipher interface {
	Encrypt(asymmetric.IPubKey, payload.IPayload) (message.IMessage, error)
	Decrypt(message.IMessage) (asymmetric.IPubKey, payload.IPayload, error)
}