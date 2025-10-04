package qb

import (
	"context"
	"time"

	"github.com/number571/go-peer/pkg/anonymity/qb/adapters"
	"github.com/number571/go-peer/pkg/anonymity/qb/queue"
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/go-peer/pkg/types"
)

type (
	IHandlerF func(context.Context, INode, uint64, asymmetric.IPubKey, []byte) ([]byte, error)
)

type INode interface {
	types.IRunner
	HandleFunc(uint32, IHandlerF) INode

	GetLogger() logger.ILogger
	GetSettings() ISettings
	GetAdapter() adapters.IAdapter
	GetKVDatabase() database.IKVDatabase
	GetMapPubKeys() asymmetric.IMapPubKeys
	GetQBProcessor() queue.IQBProblemProcessor

	WithDecryptors(...client.IDecryptor) INode

	SendPayload(context.Context, asymmetric.IPubKey, payload.IPayload64) error
	FetchPayload(context.Context, asymmetric.IPubKey, payload.IPayload32) ([]byte, error)
}

type ISettings interface {
	GetServiceName() string
	GetFetchTimeout() time.Duration
}
