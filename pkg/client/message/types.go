package message

import (
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/types"
)

type IMessage interface {
	types.IConverter
	IsValid(ISettings) bool

	GetHead() IHead
	GetBody() IBody
	GetPayload() payload.IPayload
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
	GetHash() []byte
	GetSign() []byte
	GetProof() uint64
}
