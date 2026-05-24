package qb

import (
	"context"
	"time"

	"github.com/number571/go-peer/pkg/anonymity/qb/adapters"
	"github.com/number571/go-peer/pkg/anonymity/qb/queue"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/go-peer/pkg/types"
)

type (
	IHandlerF func(context.Context, INode, layer2.IParticipantKey, []byte) ([]byte, error)
)

type INode interface {
	types.IRunner
	HandleFunc(uint32, IHandlerF) INode

	GetLogger() logger.ILogger
	GetSettings() ISettings
	GetAdapter() adapters.IAdapter
	GetKVDatabase() database.IKVDatabase
	GetKeysContainer() layer2.IKeysContainer
	GetQBProcessor() queue.IQBProblemProcessor

	SendPayload(context.Context, layer2.IParticipantKey, payload.IPayload64) error
	FetchPayload(context.Context, layer2.IParticipantKey, payload.IPayload32) ([]byte, error)
}

type ISettings interface {
	GetServiceName() string
	GetFetchTimeout() time.Duration
}
