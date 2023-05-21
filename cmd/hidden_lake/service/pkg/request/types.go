package request

type IRequest interface {
	ToBytes() []byte

	WithHead(map[string]string) IRequest
	WithBody([]byte) IRequest

	GetMethod() string
	GetHost() string
	GetPath() string
	GetHead() map[string]string
	GetBody() []byte
}
