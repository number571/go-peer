package network

import (
	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/local/client"
	"github.com/number571/go-peer/local/message"
	"github.com/number571/go-peer/local/routing"
)

type iRouterF func(INode) []asymmetric.IPubKey
type iHandlerF func(INode, message.IMessage) []byte

type INode interface {
	Client() client.IClient

	iModifier
	iHandler
	iConnect
}

type iModifier interface {
	Checker() iChecker
	Pseudo() iPseudo
	Online() iOnline
	F2F() iF2F
}

type iHandler interface {
	WithResponse(iRouterF) INode
	Request(routing.IRoute, message.IMessage) ([]byte, error)
	Handle([]byte, iHandlerF) INode
	Listen(string) error
	Close()
}

type iConnect interface {
	InConnections(string) bool
	Connections() []string
	Connect(string) error
	Disconnect(string)
}

type iOnline interface {
	iStatus
}

type iF2F interface {
	iStatus
	iListPubKey
}

type iListPubKey interface {
	InList(asymmetric.IPubKey) bool
	List() []asymmetric.IPubKey
	Append(asymmetric.IPubKey)
	Remove(asymmetric.IPubKey)
}

type iPseudo interface {
	iStatus
	Request(int) iPseudo
	Sleep() iPseudo
	PubKey() asymmetric.IPubKey
	PrivKey() asymmetric.IPrivKey
}

type iStatus interface {
	Switch(bool)
	Status() bool
}

type iChecker interface {
	ListWithInfo() []iCheckerInfo

	iStatus
	iListPubKey
}

type iCheckerInfo interface {
	Online() bool
	PubKey() asymmetric.IPubKey
}
