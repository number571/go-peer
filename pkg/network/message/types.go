package message

import (
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/types"
)

type IMessage interface {
	types.IConverter

	GetProof() uint64
	GetSalt() []byte
	GetHash() []byte
	GetPayload() payload.IPayload
}

type ISettings interface {
	GetWorkSizeBits() uint64
	GetNetworkKey() string
}
