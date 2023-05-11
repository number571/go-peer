package adapters

import "github.com/number571/go-peer/pkg/payload"

type IPayload interface {
	GetHead() uint32
	GetBody() []byte
	ToOrigin() payload.IPayload
}
