package local

import (
	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/settings"
)

type ISelector interface {
	Length() uint64
	Shuffle() ISelector
	Return(uint64) []crypto.IPubKey
}

type IRoute interface {
	WithRedirects(crypto.IPrivKey, []crypto.IPubKey) IRoute

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

type iCipher interface {
	Encrypt(IRoute, IMessage) (IMessage, []byte)
	Decrypt(IMessage) (IMessage, []byte)
}

type IClient interface {
	iKeeper
	iCipher
}

type IStorage interface {
	Write(string, string, []byte) error
	Read(string, string) ([]byte, error)
	Delete(string, string) error
}
