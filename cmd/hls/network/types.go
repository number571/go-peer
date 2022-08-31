package network

type IRequest interface {
	Bytes() []byte

	WithHead(map[string]string) IRequest
	WithBody([]byte) IRequest

	Method() string
	Host() string
	Path() string
	Head() map[string]string
	Body() []byte
}
