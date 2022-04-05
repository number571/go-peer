package network

import (
	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/local"
)

type iHandler func(INode, local.IMessage) []byte

type INode interface {
	Listen(string) error
	Close()

	Client() local.IClient
	F2F() iF2F

	Handle([]byte, iHandler) INode
	Request(local.IRoute, local.IMessage) ([]byte, error)

	Connect(string) error
	Disconnect(string)

	InConnections(string) bool
	Connections() []string
}

type iF2F interface {
	Set(bool)
	Status() bool

	Append(crypto.IPubKey)
	Remove(crypto.IPubKey)

	InList(crypto.IPubKey) bool
	List() []crypto.IPubKey
}
