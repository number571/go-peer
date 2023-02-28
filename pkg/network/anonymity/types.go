package anonymity

import (
	"time"

	"github.com/number571/go-peer/pkg/client/queue"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/go-peer/pkg/types"

	"github.com/number571/go-peer/pkg/payload"
)

type (
	IHandlerF func(INode, asymmetric.IPubKey, []byte, []byte) []byte
)

type ISettings interface {
	GetServiceName() string
	GetTimeWait() time.Duration
	GetNetworkMask() uint64
	GetRetryEnqueue() uint64
}

type INode interface {
	types.ICommand

	GetSettings() ISettings
	GetKeyValueDB() database.IKeyValueDB
	GetNetworkNode() network.INode
	GetMessageQueue() queue.IMessageQueue
	GetListPubKeys() asymmetric.IListPubKeys
	GetLogger() logger.ILogger

	HandleFunc(uint32, IHandlerF) INode
	BroadcastPayload(recv asymmetric.IPubKey, pl payload.IPayload) error
	FetchPayload(recv asymmetric.IPubKey, pl payload.IPayload) ([]byte, error)
}

type iHead interface {
	Uint64() uint64
	GetRoute() uint32
	GetAction() uint32
}
