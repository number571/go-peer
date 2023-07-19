package request

import "github.com/number571/go-peer/pkg/types"

type IRequest interface {
	types.IConverter

	WithHead(map[string]string) IRequest
	WithBody([]byte) IRequest

	GetMethod() string
	GetHost() string
	GetPath() string
	GetHead() map[string]string
	GetBody() []byte
}
