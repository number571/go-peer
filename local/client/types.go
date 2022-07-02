package client

import (
	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/local/message"
	"github.com/number571/go-peer/local/payload"
	"github.com/number571/go-peer/local/routing"
	"github.com/number571/go-peer/settings"
)

type IClient interface {
	iKeeper
	iCipher
}

type iKeeper interface {
	PubKey() asymmetric.IPubKey
	PrivKey() asymmetric.IPrivKey
	Settings() settings.ISettings
}

type iCipher interface {
	Encrypt(routing.IRoute, payload.IPayload) message.IMessage
	Decrypt(message.IMessage) (asymmetric.IPubKey, payload.IPayload)
}
