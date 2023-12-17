package anonymity

import (
	"time"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/anonymity/adapters"
	"github.com/number571/go-peer/pkg/network/anonymity/queue"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/go-peer/pkg/types"
)

type (
	IHandlerF func(INode, asymmetric.IPubKey, []byte) ([]byte, error)
)

type INode interface {
	types.IRunner

	GetSettings() ISettings
	GetWrapperDB() IWrapperDB
	GetNetworkNode() network.INode
	GetMessageQueue() queue.IMessageQueue
	GetListPubKeys() asymmetric.IListPubKeys
	GetLogger() logger.ILogger

	HandleFunc(uint32, IHandlerF) INode

	BroadcastPayload(asymmetric.IPubKey, adapters.IPayload) error
	FetchPayload(asymmetric.IPubKey, adapters.IPayload) ([]byte, error)
}

type ISettings interface {
	GetServiceName() string
	GetNetworkMask() uint64
	GetRetryEnqueue() uint64
	GetFetchTimeWait() time.Duration
}

type IWrapperDB interface {
	types.ICloser

	Get() database.IKVDatabase
	Set(database.IKVDatabase) IWrapperDB
}
