package anonymity

import (
	"context"
	"time"

	"github.com/number571/go-peer/pkg/anonymity/queue"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/go-peer/pkg/types"

	net_message "github.com/number571/go-peer/pkg/network/message"
)

type IProducerF func(context.Context, net_message.IMessage) error
type IConsumerF func(context.Context) (net_message.IMessage, error)

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
	GetMessageQueue() queue.IQBProblemProcessor
	GetMapPubKeys() asymmetric.IMapPubKeys

	SendPayload(context.Context, asymmetric.IPubKey, payload.IPayload64) error
	FetchPayload(context.Context, asymmetric.IPubKey, payload.IPayload32) ([]byte, error)
}

type ISettings interface {
	GetServiceName() string
	GetFetchTimeout() time.Duration
}
