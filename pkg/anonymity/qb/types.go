package qb

import (
	"context"
	"time"

	"github.com/number571/go-peer/pkg/anonymity/qb/adapters"
	"github.com/number571/go-peer/pkg/anonymity/qb/queue"
	"github.com/number571/go-peer/pkg/crypto/hybrid"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/go-peer/pkg/types"
)

type (
	IHandlerF func(context.Context, INode, hybrid.IParticipantKey, []byte) ([]byte, error)
)

type INode interface {
	types.IRunner
	HandleFunc(uint32, IHandlerF) INode

	GetLogger() logger.ILogger
	GetSettings() ISettings
	GetAdapter() adapters.IAdapter
	GetKVDatabase() database.IKVDatabase
	GetKeysContainer() hybrid.IKeysContainer
	GetQBProcessor() queue.IQBProblemProcessor

	SendPayload(context.Context, hybrid.IParticipantKey, payload.IPayload64) error
	FetchPayload(context.Context, hybrid.IParticipantKey, payload.IPayload32) ([]byte, error)
}

type ISettings interface {
	GetServiceName() string
	GetFetchTimeout() time.Duration
}
