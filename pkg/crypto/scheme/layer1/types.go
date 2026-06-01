package layer1

import (
	"github.com/number571/go-peer/pkg/types"
)

type IMessage interface {
	types.IConverter

	// body = payload
	GetBody() []byte

	// hash = H(payload)
	GetHash() []byte

	// hmac = HMAC(network_key, payload)
	GetHmac() []byte

	// proof = PoW(hmac)
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
