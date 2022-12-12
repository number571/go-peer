package anonymity

import (
	"time"

	"github.com/number571/go-peer/modules"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/friends"
	"github.com/number571/go-peer/modules/network"
	"github.com/number571/go-peer/modules/queue"
	"github.com/number571/go-peer/modules/storage/database"

	payload_adapter "github.com/number571/go-peer/modules/network/anonymity/adapters/payload"
)

type IHandlerF func(INode, asymmetric.IPubKey, []byte) []byte

type INode interface {
	modules.IApp

	Settings() ISettings
	KeyValueDB() database.IKeyValueDB
	Network() network.INode
	Queue() queue.IQueue
	F2F() friends.IF2F

	Handle(uint32, IHandlerF) INode
	Broadcast(recv asymmetric.IPubKey, pl payload_adapter.IPayload) error
	Request(recv asymmetric.IPubKey, pl payload_adapter.IPayload) ([]byte, error)
}

type ISettings interface {
	GetTimeWait() time.Duration
	GetNetworkMask() uint64
	GetRetryEnqueue() uint64
}

type iHead interface {
	Uint64() uint64
	GetRoute() uint32
	GetAction() uint32
}
