package message

import (
	"time"

	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/types"
)

const (
	CWithoutTimestamp = time.Duration(0)
)

type IMessage interface {
	types.IConverter

	GetHash() []byte
	GetRand() []byte
	GetTime() uint64
	GetProof() uint64

	// payload = head(32bit) || body(Nbit)
	GetPayload() payload.IPayload32
}

type IConstructSettings interface {
	ISettings
	GetParallel() uint64
	GetRandMessageSizeBytes() uint64
}

type ISettings interface {
	GetWorkSizeBits() uint64
	GetNetworkKey() string
}
