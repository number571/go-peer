package response

type IResponse interface {
	ToBytes() []byte

	GetCode() int
	GetBody() []byte
	// GetHead() map[string]string
}
