package network

import (
	"context"
	"time"

	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/types"

	net_message "github.com/number571/go-peer/pkg/network/message"
)

type IHandlerF func(
	context.Context,
	INode,
	conn.IConn,
	net_message.IMessage,
) error

type INode interface {
	types.ICloser

	Listen(context.Context) error
	HandleFunc(uint32, IHandlerF) INode

	SetVSettings(conn.IVSettings)
	GetVSettings() conn.IVSettings

	GetSettings() ISettings
	GetConnections() map[string]conn.IConn

	AddConnection(context.Context, string) error
	DelConnection(string) error

	BroadcastMessage(context.Context, net_message.IMessage) error
}

type ISettings interface {
	GetAddress() string
	GetMaxConnects() uint64
	GetReadTimeout() time.Duration
	GetWriteTimeout() time.Duration
	GetConnSettings() conn.ISettings
}
