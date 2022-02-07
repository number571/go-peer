package local

import (
	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/settings"
)

type (
	Message = *MessageT
	Route   = *RouteT
	Session = []byte
)

type Keeper interface {
	PubKey() crypto.PubKey
	PrivKey() crypto.PrivKey
	Settings() settings.Settings
}

type Cipher interface {
	Encrypt(Route, Message) (Message, Session)
	Decrypt(Message) Message
}

type Client interface {
	Keeper
	Cipher
}

type Storage interface {
	Write([]byte, string, string) error
	Read(string, string) ([]byte, error)
	Delete(string, string) error
}
