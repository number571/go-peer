package hlc

import (
	"github.com/number571/go-peer/modules/crypto/asymmetric"

	hls_network "github.com/number571/go-peer/cmd/hls/network"
	hls_settings "github.com/number571/go-peer/cmd/hls/settings"
)

// TODO: need update
type IClient interface {
	PubKey() (asymmetric.IPubKey, error)
	Online() ([]string, error)
	Friends() ([]asymmetric.IPubKey, error)
	Request(asymmetric.IPubKey, hls_network.IRequest) ([]byte, error)
}

type IBuilder interface {
	Request(asymmetric.IPubKey, hls_network.IRequest) *hls_settings.SPush
}

type IRequester interface {
	PubKey() (asymmetric.IPubKey, error)
	Online() ([]string, error)
	Friends() ([]asymmetric.IPubKey, error)
	Request(*hls_settings.SPush) ([]byte, error)
}
