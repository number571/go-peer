package netanon

import (
	"time"

	"github.com/number571/go-peer/client"
	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/friends"
	"github.com/number571/go-peer/message"
	"github.com/number571/go-peer/network"
	"github.com/number571/go-peer/payload"
	"github.com/number571/go-peer/queue"
)

type IRouterF func() []asymmetric.IPubKey
type IHandlerF func(INode, asymmetric.IPubKey, payload.IPayload) []byte

type INode interface {
	Settings() ISettings
	Close() error

	Client() client.IClient
	Network() network.INode
	Queue() queue.IQueue
	F2F() friends.IF2F

	Handle(uint32, IHandlerF) INode
	Broadcast(message.IMessage) error
	Request(recv asymmetric.IPubKey, pl payload.IPayload) ([]byte, error)
}

type ISettings interface {
	GetResponsePeriod() uint64
	GetRetryEnqueue() uint64
	GetTimeWait() time.Duration
}

type iHead interface {
	Uint64() uint64
	Routes() uint32
	Actions() uint32
}
