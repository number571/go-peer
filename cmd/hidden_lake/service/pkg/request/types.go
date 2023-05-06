package request

type IRequest interface {
	ToBytes() []byte

	WithHead(map[string]string) IRequest
	WithBody([]byte) IRequest

	Method() string
	Host() string
	Path() string
	Head() map[string]string
	Body() []byte
}
