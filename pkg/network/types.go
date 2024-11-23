package network

import (
	"context"
	"io"
	"time"

	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/storage/cache"

	net_message "github.com/number571/go-peer/pkg/network/message"
)

type IHandlerF func(
	context.Context,
	INode,
	conn.IConn,
	net_message.IMessage,
) error

type INode interface {
	io.Closer

	Listen(context.Context) error
	HandleFunc(uint32, IHandlerF) INode

	GetSettings() ISettings
	GetCacheSetter() cache.ICacheSetter

	GetConnections() map[string]conn.IConn
	AddConnection(context.Context, string) error
	DelConnection(string) error

	BroadcastMessage(context.Context, net_message.IMessage) error
}

type ISettings interface {
	GetConnSettings() conn.ISettings
	GetAddress() string
	GetMaxConnects() uint64
	GetReadTimeout() time.Duration
	GetWriteTimeout() time.Duration
}
