package network

import (
	"net"
	"time"

	"github.com/number571/go-peer/modules"
	"github.com/number571/go-peer/modules/payload"
)

type IHandlerF func(INode, IConn, payload.IPayload)

type INode interface {
	Settings() ISettings
	Connections() map[string]IConn

	Handle(uint64, IHandlerF) INode
	Broadcast(payload.IPayload) error

	Listen(string) error
	modules.ICloser

	Connect(string) IConn
	Disconnect(string) error
}

type ISettings interface {
	GetNetworkKey() string
	GetCapacity() uint64
	GetMessageSize() uint64
	GetMaxConnects() uint64
	GetTimeWait() time.Duration
}

type IConn interface {
	Socket() net.Conn
	Request(IMessage) IMessage

	Write(IMessage) error
	Read() IMessage
	modules.ICloser
}

type IMessage interface {
	Hash() []byte
	Payload() payload.IPayload
	Bytes() []byte
}

type iPackage interface {
	SizeToBytes() []byte
	BytesToSize() uint64
}
