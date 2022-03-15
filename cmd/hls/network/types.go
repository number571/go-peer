package network

type IRequest interface {
	ToBytes() []byte

	WithHead(map[string]string) IRequest
	WithBody([]byte) IRequest

	Host() string
	Path() string
	Method() string
	Head() map[string]string
	Body() []byte
}
