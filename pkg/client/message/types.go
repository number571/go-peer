package message

import "github.com/number571/go-peer/pkg/payload"

type IMessage interface {
	GetHead() IHead
	GetBody() IBody

	IsValid(IParams) bool
	ToBytes() []byte
}

type IParams interface {
	GetMessageSize() uint64
	GetWorkSize() uint64
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
