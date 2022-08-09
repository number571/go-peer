package message

import "github.com/number571/go-peer/payload"

type IMessage interface {
	Head() iHead
	Body() iBody
	Bytes() []byte
}

type iHead interface {
	Sender() []byte
	Session() []byte
	Salt() []byte
}

type iBody interface {
	Payload() payload.IPayload
	Hash() []byte
	Sign() []byte
	Proof() uint64
}
