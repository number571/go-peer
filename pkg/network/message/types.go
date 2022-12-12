package message

import "github.com/number571/go-peer/pkg/payload"

type IMessage interface {
	Hash() []byte
	Payload() payload.IPayload
	Bytes() []byte
}
