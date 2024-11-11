package message

import (
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/types"
)

type IMessage interface {
	types.IConverter

	GetHash() []byte
	GetProof() uint64

	// payload = head(32bit) || body(Nbit)
	GetPayload() payload.IPayload32
}

type IConstructSettings interface {
	GetSettings() ISettings
	GetParallel() uint64
}

type ISettings interface {
	GetWorkSizeBits() uint64
	GetNetworkKey() string
}
