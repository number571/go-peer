package netanon

import (
	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/local/client"
	"github.com/number571/go-peer/local/message"
	"github.com/number571/go-peer/local/payload"
	"github.com/number571/go-peer/network"
)

type IRouterF func() []asymmetric.IPubKey
type IHandlerF func(INode, asymmetric.IPubKey, payload.IPayload) []byte

type INode interface {
	Client() client.IClient
	Network() network.INode

	Handle(uint32, IHandlerF) INode
	WithRouter(IRouterF) INode

	Broadcast(message.IMessage) error
	Request(recv asymmetric.IPubKey, pl payload.IPayload) ([]byte, error)

	F2F() iF2F
	Online() iOnline
	Checker() iChecker
	Pseudo() iPseudo

	Close() error
}

type iOnline interface {
	iStatus
}

type iF2F interface {
	iStatus
	iListPubKey
}

type iHead interface {
	Uint64() uint64
	Routes() uint32
	Actions() uint32
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

type iListPubKey interface {
	InList(asymmetric.IPubKey) bool
	List() []asymmetric.IPubKey
	Append(asymmetric.IPubKey)
	Remove(asymmetric.IPubKey)
}

type iPseudo interface {
	iStatus

	request(int) iPseudo
	sleep() iPseudo
	privKey() asymmetric.IPrivKey
}
