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

func (builder *sBuilder) Connect(connect string) *hls_settings.SConnect {
	return &hls_settings.SConnect{
		FConnect: connect,
	}
}

func (builder *sBuilder) Friend(aliasName string, pubKey asymmetric.IPubKey) *hls_settings.SFriend {
	if pubKey == nil {
		return &hls_settings.SFriend{
			FAliasName: aliasName,
		}
	}
	return &hls_settings.SFriend{
		FAliasName: aliasName,
		FPublicKey: pubKey.String(),
	}
}

func (builder *sBuilder) Push(recv asymmetric.IPubKey, req hls_network.IRequest) *hls_settings.SPush {
	return &hls_settings.SPush{
		FReceiver: recv.String(),
		FHexData:  encoding.HexEncode(req.Bytes()),
	}
}
