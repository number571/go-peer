package anonymity

import (
	"context"
	"time"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/database"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/anonymity/queue"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/types"
)

type IHandlerF func(
	context.Context,
	INode,
	asymmetric.IPubKey,
	[]byte,
) ([]byte, error)

type INode interface {
	types.IRunner
	HandleFunc(uint32, IHandlerF) INode

	GetLogger() logger.ILogger
	GetSettings() ISettings
	GetKVDatabase() database.IKVDatabase
	GetNetworkNode() network.INode
	GetMessageQueue() queue.IMessageQueue
	GetListPubKeys() asymmetric.IListPubKeys

	SendPayload(context.Context, asymmetric.IPubKey, payload.IPayload64) error
	FetchPayload(context.Context, asymmetric.IPubKey, payload.IPayload32) ([]byte, error)
}

type ISettings interface {
	GetServiceName() string
	GetF2FDisabled() bool
	GetNetworkMask() uint64
	GetRetryEnqueue() uint64
	GetFetchTimeout() time.Duration
}
