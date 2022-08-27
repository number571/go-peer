package anonymity

import (
	"time"

	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/friends"
	"github.com/number571/go-peer/message"
	"github.com/number571/go-peer/network"
	"github.com/number571/go-peer/payload"
	"github.com/number571/go-peer/queue"
	"github.com/number571/go-peer/storage/database"

	payload_adapter "github.com/number571/go-peer/network/anonymity/adapters/payload"
)

type IHandlerF func(INode, asymmetric.IPubKey, payload.IPayload) []byte

type INode interface {
	Settings() ISettings
	KeyValueDB() database.IKeyValueDB
	Network() network.INode
	Queue() queue.IQueue
	F2F() friends.IF2F

	Handle(uint32, IHandlerF) INode
	Broadcast(message.IMessage) error
	Request(recv asymmetric.IPubKey, pl payload_adapter.IPayload) ([]byte, error)

	Run() error
	Close() error
}

type ISettings interface {
	GetTimeWait() time.Duration
	GetRetryEnqueue() uint64
}

type iHead interface {
	Uint64() uint64
	Routes() uint32
	Actions() uint32
}
