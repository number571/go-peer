package anonymity

import (
	"context"
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
	IHandlerF func(context.Context, INode, asymmetric.IPubKey, []byte) ([]byte, error)
)

type INode interface {
	types.IRunner
	HandleFunc(uint32, IHandlerF) INode

	GetSettings() ISettings
	GetDBWrapper() IDBWrapper
	GetNetworkNode() network.INode
	GetMessageQueue() queue.IMessageQueue
	GetListPubKeys() asymmetric.IListPubKeys
	GetLogger() logger.ILogger

	BroadcastPayload(context.Context, asymmetric.IPubKey, adapters.IPayload) error
	FetchPayload(context.Context, asymmetric.IPubKey, adapters.IPayload) ([]byte, error)
}

type ISettings interface {
	GetServiceName() string
	GetNetworkMask() uint64
	GetRetryEnqueue() uint64
	GetFetchTimeWait() time.Duration
}

type IDBWrapper interface {
	types.ICloser

	Get() database.IKVDatabase
	Set(database.IKVDatabase) IDBWrapper
}
