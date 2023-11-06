package message

import "github.com/number571/go-peer/pkg/payload"

type IMessage interface {
	ToBytes() []byte

	GetProof() uint64
	GetHash() []byte
	GetPayload() payload.IPayload
}

type ISettings interface {
	GetWorkSizeBits() uint64
	GetNetworkKey() string
}
