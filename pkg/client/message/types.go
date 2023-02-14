package message

import "github.com/number571/go-peer/pkg/payload"

type IMessage interface {
	Head() IHead
	Body() IBody

	IsValid(IParams) bool
	Bytes() []byte
}

type IParams interface {
	GetMessageSize() uint64
	GetWorkSize() uint64
}

type IHead interface {
	Sender() []byte
	Session() []byte
	Salt() []byte
}

type IBody interface {
	Payload() payload.IPayload
	Hash() []byte
	Sign() []byte
	Proof() uint64
}
