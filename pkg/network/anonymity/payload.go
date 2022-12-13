package anonymity

import "github.com/number571/go-peer/pkg/payload"

func NewPayload(head uint32, body []byte) payload.IPayload {
	return payload.NewPayload(uint64(head), body)
}
