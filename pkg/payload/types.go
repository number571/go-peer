package payload

type IPayload64 interface {
	GetHead() uint64
	iPayload
}

type IPayload32 interface {
	GetHead() uint32
	iPayload
}

type iPayload interface {
	GetBody() []byte
	ToBytes() []byte
}
