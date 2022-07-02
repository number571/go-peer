package message

type IMessage interface {
	Head() iHead
	Body() iBody
	Bytes() []byte
}

type IPayload interface {
	Head() uint64
	Body() []byte
	Bytes() []byte
}

type iHead interface {
	Sender() []byte
	Session() []byte
	Salt() []byte
}

type iBody interface {
	Payload() IPayload
	Hash() []byte
	Sign() []byte
	Proof() uint64
}
