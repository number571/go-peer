package client

import (
	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/local/message"
	"github.com/number571/go-peer/local/routing"
	"github.com/number571/go-peer/settings"
)

type iKeeper interface {
	PubKey() asymmetric.IPubKey
	PrivKey() asymmetric.IPrivKey
	Settings() settings.ISettings
}

type iCipher interface {
	Encrypt(routing.IRoute, message.IMessage) (message.IMessage, []byte)
	Decrypt(message.IMessage) (message.IMessage, []byte)
}

type IClient interface {
	iKeeper
	iCipher
}
