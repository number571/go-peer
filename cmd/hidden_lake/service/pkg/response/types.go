package response

import "github.com/number571/go-peer/pkg/types"

type IResponse interface {
	types.IConverter

	WithHead(map[string]string) IResponse
	WithBody(pBody []byte) IResponse

	GetCode() int
	GetHead() map[string]string
	GetBody() []byte
}
