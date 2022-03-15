package network

import (
	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/local"
)

type iF2F interface {
	Status() bool
	Switch()

	Append(crypto.IPubKey)
	Remove(crypto.IPubKey)

	InList(crypto.IPubKey) bool
	List() []crypto.IPubKey
}

type (
	Title    = []byte
	Response = []byte
	Address  = string
	Handler  = func(local.IClient, local.IMessage) []byte
)
type INode interface {
	Listen(Address) error
	Close()

	Client() local.IClient
	F2F() iF2F

	Handle(Title, Handler) INode
	Request(local.IRoute, local.IMessage) (Response, error)

	Connect(Address) error
	Disconnect(Address)

	InConnections(Address) bool
	Connections() []Address
}
