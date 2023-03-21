package anonymity

import "github.com/number571/go-peer/pkg/payload"

func NewPayload(pHead uint32, pBody []byte) payload.IPayload {
	return payload.NewPayload(uint64(pHead), pBody)
}
