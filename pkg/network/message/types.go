package message

import "github.com/number571/go-peer/pkg/payload"

type IMessage interface {
	GetHash() []byte
	GetPayload() payload.IPayload
	GetBytes() []byte
}
