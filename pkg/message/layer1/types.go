package layer1

import (
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/types"
)

type IMessage interface {
	types.IConverter

	// payload = head(32bit) || body(Nbit)
	GetPayload() payload.IPayload32

	// hash = H(payload)
	GetHash() []byte

	// proof = PoW(HMAC(network_key, hash))
	GetProof() uint64
}

type IConstructSettings interface {
	GetSettings() ISettings
	GetParallel() uint64
}

type ISettings interface {
	GetWorkSizeBits() uint64
	GetNetworkKey() string
}
