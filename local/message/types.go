package message

type IMessage interface {
	Head() iHead
	Body() iBody

	ToPackage() IPackage
}

type IPackage interface {
	Size() uint64
	Bytes() []byte

	SizeToBytes() []byte
	BytesToSize() uint64

	ToMessage() IMessage
}

type iHead interface {
	Sender() []byte
	Session() []byte
	Salt() []byte
}

type iBody interface {
	Data() []byte
	Hash() []byte
	Sign() []byte
	Proof() uint64
}
