package network

import (
	"time"

	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/types"
)

type IHandlerF func(INode, conn.IConn, []byte)

type INode interface {
	types.ICommand

	GetSettings() ISettings
	GetConnections() map[string]conn.IConn

	AddConnect(string) error
	DelConnect(string) error

	HandleFunc(uint64, IHandlerF) INode
	BroadcastPayload(payload.IPayload) error
}

type ISettings interface {
	GetAddress() string
	GetCapacity() uint64
	GetMaxConnects() uint64
	GetWriteTimeout() time.Duration
	GetConnSettings() conn.ISettings
}
