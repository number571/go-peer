package network

import (
	"context"
	"time"

	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/types"
)

type (
	IHandlerF func(context.Context, INode, conn.IConn, message.IMessage) error
)

type INode interface {
	types.ICloser
	Listen(context.Context) error

	GetSettings() ISettings
	GetConnections() map[string]conn.IConn

	AddConnection(context.Context, string) error
	DelConnection(string) error

	HandleFunc(uint64, IHandlerF) INode
	BroadcastMessage(context.Context, message.IMessage) error
}

type ISettings interface {
	GetAddress() string
	GetMaxConnects() uint64
	GetReadTimeout() time.Duration
	GetWriteTimeout() time.Duration
	GetConnSettings() conn.ISettings
}
