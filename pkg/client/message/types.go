package message

import (
	"github.com/number571/go-peer/pkg/types"
)

type IMessage interface {
	types.IConverter

	GetEnck() []byte
	GetEncd() []byte
}
