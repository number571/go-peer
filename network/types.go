package network

import (
	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/local"
)

type (
	Title    = []byte
	Response = []byte
	Address  = string
	Handler  = func(local.Client, local.Message) []byte
)
type Node interface {
	Listen(Address) error
	Close()

	Client() local.Client
	F2F() F2F

	Handle(Title, Handler) Node
	Request(local.Route, local.Message) (Response, error)

	Connect(Address) error
	Disconnect(Address)

	InConnections(Address) bool
	Connections() []Address
}

type F2F interface {
	Status() bool
	Switch()

	Append(crypto.PubKey)
	Remove(crypto.PubKey)

	InList(crypto.PubKey) bool
	List() []crypto.PubKey
}
