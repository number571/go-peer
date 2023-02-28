package payload

type IPayload interface {
	GetHead() uint64
	GetBody() []byte
	ToBytes() []byte
}
