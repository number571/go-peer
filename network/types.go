package network

import (
	"net"
	"time"

	"github.com/number571/go-peer/payload"
)

type IHandlerF func(INode, IConn, payload.IPayload)

type INode interface {
	Settings() ISettings

	Handle(uint64, IHandlerF) INode
	Broadcast(payload.IPayload) error

	Listen(string) error
	Close() error

	Connect(string) IConn
	Disconnect(IConn) error

	Connections() []IConn
}

type ISettings interface {
	GetRetryNum() uint64
	GetCapacity() uint64
	GetPackageSize() uint64
	GetMaxConnects() uint64
	GetMaxMessages() uint64
	GetTimeWait() time.Duration
}

type IConn interface {
	Socket() net.Conn
	Request(IMessage) IMessage

	Write(IMessage) error
	Read() IMessage
	Close() error
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
