package network

import (
	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/offline/client"
	"github.com/number571/go-peer/offline/message"
	"github.com/number571/go-peer/offline/routing"
)

type iRouter func(INode) []asymmetric.IPubKey
type iHandler func(INode, message.IMessage) []byte

type INode interface {
	WithResponseRouter(iRouter) INode

	Request(routing.IRoute, message.IMessage) ([]byte, error)
	Handle([]byte, iHandler) INode
	Listen(string) error
	Close()

	InConnections(string) bool
	Connections() []string
	Connect(string) error
	Disconnect(string)

	Client() client.IClient
	Checker() iChecker
	Pseudo() iPseudo
	Online() iOnline
	F2F() iF2F
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
