package hlc

import (
	hls_network "github.com/number571/go-peer/cmd/hls/network"
	hls_settings "github.com/number571/go-peer/cmd/hls/settings"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/encoding"
)

var (
	_ IBuilder = &sBuilder{}
)

type sBuilder struct {
}

func NewBuilder() IBuilder {
	return &sBuilder{}
}

func (builder *sBuilder) Request(recv asymmetric.IPubKey, req hls_network.IRequest) *hls_settings.SRequest {
	return &hls_settings.SRequest{
		FReceiver: recv.String(),
		FHexData:  encoding.HexEncode(req.Bytes()),
	}
}
