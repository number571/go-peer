package message

import "github.com/number571/go-peer/pkg/payload"

type IMessage interface {
	GetHead() IHead
	GetBody() IBody

	IsValid(ISettings) bool
	ToBytes() []byte
}

type ISettings interface {
	GetMessageSizeBytes() uint64
	GetWorkSizeBits() uint64
}

type IHead interface {
	GetSender() []byte
	GetSession() []byte
	GetSalt() []byte
}

type IBody interface {
	GetPayload() payload.IPayload
	GetHash() []byte
	GetSign() []byte
	GetProof() uint64
}
