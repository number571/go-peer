package layer2

import (
	"github.com/number571/go-peer/pkg/types"
)

type IMessage interface {
	types.IConverter

	GetEnck() []byte
	GetEncd() []byte
}
