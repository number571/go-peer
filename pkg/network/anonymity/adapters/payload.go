package adapters

import (
	"github.com/number571/go-peer/pkg/payload"
)

var (
	_ IPayload = &sPayload{}
)

type sPayload struct {
	fPayload payload.IPayload
}

func NewPayload(pHead uint32, pBody []byte) IPayload {
	return &sPayload{
		fPayload: payload.NewPayload(uint64(pHead), pBody),
	}
}

func (p *sPayload) GetHead() uint32 {
	return uint32(p.fPayload.GetHead())
}

func (p *sPayload) GetBody() []byte {
	return p.fPayload.GetBody()
}

func (p *sPayload) ToOrigin() payload.IPayload {
	return p.fPayload
}
