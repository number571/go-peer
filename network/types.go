package network

import (
	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/local"
)

type F2F interface {
	State() bool
	Switch()

	InList(crypto.PubKey) bool
	List() []crypto.PubKey

	Append(crypto.PubKey)
	Remove(crypto.PubKey)
}

type (
	Handler = func(local.Client, local.Message) []byte
)
type Node interface {
	Client() local.Client
	F2F() F2F

	Listen(string) error
	Close()

	Handle([]byte, Handler) Node
	Broadcast(local.Route, local.Message) ([]byte, error)

	InConnections(string) bool
	Connections() []string

	Connect(string) error
	Disconnect(string)
}
