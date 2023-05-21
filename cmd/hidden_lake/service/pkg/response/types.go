package response

type IResponse interface {
	ToBytes() []byte

	WithBody(pBody []byte) IResponse
	WithHead(map[string]string) IResponse

	GetCode() int
	GetBody() []byte
	GetHead() map[string]string
}
