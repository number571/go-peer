package client

import (
	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/message"
	"github.com/number571/go-peer/payload"
	"github.com/number571/go-peer/routing"
)

type IClient interface {
	Settings() ISettings

	iKeeper
	iCipher
}

type ISettings interface {
	GetRandomSize() uint64
	GetWorkSize() uint64
}

type iKeeper interface {
	PubKey() asymmetric.IPubKey
	PrivKey() asymmetric.IPrivKey
}

type iCipher interface {
	Encrypt(routing.IRoute, payload.IPayload) message.IMessage
	Decrypt(message.IMessage) (asymmetric.IPubKey, payload.IPayload)
}
