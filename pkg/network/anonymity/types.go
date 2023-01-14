package anonymity

import (
	"time"

	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/client/queue"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/friends"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/go-peer/pkg/types"

	"github.com/number571/go-peer/pkg/payload"
)

type (
	IHandlerF  func(INode, asymmetric.IPubKey, []byte) []byte
	IDownloadF func() <-chan message.IMessage
	IUploadF   func(message.IMessage) error
)

type INode interface {
	types.IApp

	Settings() ISettings
	KeyValueDB() database.IKeyValueDB
	Network() network.INode
	Queue() queue.IQueue
	F2F() friends.IF2F

	Handle(uint32, IHandlerF) INode
	Broadcast(recv asymmetric.IPubKey, pl payload.IPayload) error
	Request(recv asymmetric.IPubKey, pl payload.IPayload) ([]byte, error)
}

type ISettings interface {
	GetTimeWait() time.Duration
	GetNetworkMask() uint64
	GetRetryEnqueue() uint64
	GetTraffic() ITraffic
}

type ITraffic interface {
	Download() IDownloadF
	Upload() IUploadF
}

type iHead interface {
	Uint64() uint64
	GetRoute() uint32
	GetAction() uint32
}
