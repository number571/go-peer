package network

import (
	"context"
	"time"

	"github.com/number571/go-peer/pkg/message/layer1"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/go-peer/pkg/types"
)

type IHandlerF func(
	context.Context,
	INode,
	conn.IConn,
	layer1.IMessage,
) error

type INode interface {
	types.IRunner
	HandleFunc(uint32, IHandlerF) INode

	GetSettings() ISettings
	GetCacheSetter() cache.ICacheSetter

	GetConnections() map[string]conn.IConn
	AddConnection(context.Context, string) error
	DelConnection(string) error

	BroadcastMessage(context.Context, layer1.IMessage) error
}

type ISettings interface {
	GetConnSettings() conn.ISettings
	GetAddress() string
	GetMaxConnects() uint64
	GetReadTimeout() time.Duration
	GetWriteTimeout() time.Duration
}
