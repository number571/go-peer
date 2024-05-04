package message

import (
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/types"
)

type IMessage interface {
	types.IConverter

	GetHash() []byte
	GetVoid() []byte
	GetProof() uint64
	GetPayload() payload.IPayload64
}

type IConstructSettings interface {
	ISettings
	GetParallel() uint64
	GetLimitVoidSizeBytes() uint64
}

type ISettings interface {
	GetWorkSizeBits() uint64
	GetNetworkKey() string
}
