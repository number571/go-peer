package local

import (
	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/settings"
)

type IRoute interface {
	Receiver() crypto.IPubKey
	PSender() crypto.IPrivKey
	List() []crypto.IPubKey
}

type iHead interface {
	Sender() []byte
	Session() []byte
	Salt() []byte
}

type iBody interface {
	Data() []byte
	Hash() []byte
	Sign() []byte
	Proof() uint64
}

type IMessage interface {
	Head() iHead
	Body() iBody

	ToPackage() IPackage
}

type IPackage interface {
	Size() uint64
	Bytes() []byte

	SizeToBytes() []byte
	BytesToSize() uint64

	ToMessage() IMessage
}

type iKeeper interface {
	PubKey() crypto.IPubKey
	PrivKey() crypto.IPrivKey
	Settings() settings.ISettings
}

type (
	Session = []byte
	Title   = []byte
)
type iCipher interface {
	Encrypt(IRoute, IMessage) (IMessage, Session)
	Decrypt(IMessage) (IMessage, Title)
}

type IClient interface {
	iKeeper
	iCipher
}

type (
	Identifier = string
	Password   = string
)
type IStorage interface {
	Write(Identifier, Password, []byte) error
	Read(Identifier, Password) ([]byte, error)
	Delete(Identifier, Password) error
}
