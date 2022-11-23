package network

import (
	"github.com/number571/go-peer/modules"
	"github.com/number571/go-peer/modules/network/conn"
	"github.com/number571/go-peer/modules/payload"
)

type IHandlerF func(INode, conn.IConn, payload.IPayload)

type INode interface {
	Settings() ISettings
	Connections() map[string]conn.IConn

	Handle(uint64, IHandlerF) INode
	Broadcast(payload.IPayload) error

	Listen(string) error
	modules.ICloser

	Connect(string) (conn.IConn, error)
	Disconnect(string) error
}

type ISettings interface {
	GetCapacity() uint64
	GetMaxConnects() uint64
	GetConnSettings() conn.ISettings
}
