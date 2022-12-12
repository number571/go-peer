package network

import (
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/types"
)

type IHandlerF func(INode, conn.IConn, []byte)

type INode interface {
	Settings() ISettings
	Connections() map[string]conn.IConn

	Handle(uint64, IHandlerF) INode
	Broadcast(payload.IPayload) error

	Listen(string) error
	types.ICloser

	Connect(string) (conn.IConn, error)
	Disconnect(string) error
}

type ISettings interface {
	GetCapacity() uint64
	GetMaxConnects() uint64
	GetConnSettings() conn.ISettings
}
