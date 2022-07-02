package payload

type IPayload interface {
	Head() uint64
	Body() []byte
	Bytes() []byte
}
