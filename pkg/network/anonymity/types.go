package anonymity

import (
	"time"

	"github.com/number571/go-peer/pkg/client/message"
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
	GetWrapperDB() IWrapperDB
	GetNetworkNode() network.INode
	GetMessageQueue() queue.IMessageQueue
	GetListPubKeys() asymmetric.IListPubKeys
	GetLogger() logger.ILogger

	HandleFunc(uint32, IHandlerF) INode
	HandleMessage(message.IMessage) // in runtime

	BroadcastPayload(IFormatType, asymmetric.IPubKey, payload.IPayload) error
	FetchPayload(asymmetric.IPubKey, payload.IPayload) ([]byte, error)
}

type iHead interface {
	Uint64() uint64
	GetRoute() uint32
	GetAction() uint32
}

type IWrapperDB interface {
	types.ICloser

	Get() database.IKeyValueDB
	Set(database.IKeyValueDB) IWrapperDB
}
