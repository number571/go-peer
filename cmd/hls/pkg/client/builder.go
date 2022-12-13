package client

import (
	"github.com/number571/go-peer/cmd/hls/pkg/request"
	pkg_settings "github.com/number571/go-peer/cmd/hls/pkg/settings"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
)

var (
	_ IBuilder = &sBuilder{}
)

type sBuilder struct {
}

func NewBuilder() IBuilder {
	return &sBuilder{}
}

func (builder *sBuilder) PrivKey(privKey asymmetric.IPrivKey) *pkg_settings.SPrivKey {
	return &pkg_settings.SPrivKey{
		FPrivKey: privKey.String(),
	}
}

func (builder *sBuilder) Connect(connect string) *pkg_settings.SConnect {
	return &pkg_settings.SConnect{
		FConnect: connect,
	}
}

func (builder *sBuilder) Friend(aliasName string, pubKey asymmetric.IPubKey) *pkg_settings.SFriend {
	if pubKey == nil {
		return &pkg_settings.SFriend{
			FAliasName: aliasName,
		}
	}
	return &pkg_settings.SFriend{
		FAliasName: aliasName,
		FPublicKey: pubKey.String(),
	}
}

func (builder *sBuilder) Push(recv asymmetric.IPubKey, req request.IRequest) *pkg_settings.SPush {
	return &pkg_settings.SPush{
		FReceiver: recv.String(),
		FHexData:  encoding.HexEncode(req.Bytes()),
	}
}
