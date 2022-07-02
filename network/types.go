package network

import (
	"net"

	"github.com/number571/go-peer/local/payload"
)

type IHandlerF func(INode, IConn, payload.IPayload)

type INode interface {
	Handle(uint64, IHandlerF) INode
	Broadcast(payload.IPayload) error

	Listen(string) error
	Close() error

	Connect(string) IConn
	Disconnect(IConn) error

	Connections() []IConn
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
	Bytes() []byte
	Payload() payload.IPayload
}

type iPackage interface {
	SizeToBytes() []byte
	BytesToSize() uint64
}
