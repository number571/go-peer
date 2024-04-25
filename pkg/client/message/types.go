package message

import (
	"github.com/number571/go-peer/pkg/types"
)

type IMessage interface {
	types.IConverter
	IsValid(ISettings) bool

	GetEnck() []byte
	GetEncd() []byte
}

type ISettings interface {
	GetKeySizeBits() uint64
	GetMessageSizeBytes() uint64
}
