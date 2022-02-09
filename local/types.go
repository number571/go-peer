package local

import (
	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/settings"
)

type Keeper interface {
	PubKey() crypto.PubKey
	PrivKey() crypto.PrivKey
	Settings() settings.Settings
}

type (
	Session = []byte
	Message = *MessageT
	Route   = *RouteT
)
type Cipher interface {
	Encrypt(Route, Message) (Message, Session)
	Decrypt(Message) Message
}

type Client interface {
	Keeper
	Cipher
}

type (
	Identifier = string
	Password   = string
)
type Storage interface {
	Write(Identifier, Password, []byte) error
	Read(Identifier, Password) ([]byte, error)
	Delete(Identifier, Password) error
}
