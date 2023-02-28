package client

import (
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
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

func (builder *sBuilder) SetPrivKey(privKey asymmetric.IPrivKey) *pkg_settings.SPrivKey {
	return &pkg_settings.SPrivKey{
		FPrivKey: privKey.ToString(),
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
		FPublicKey: pubKey.ToString(),
	}
}

func (builder *sBuilder) Request(recv asymmetric.IPubKey, req request.IRequest) *pkg_settings.SRequest {
	return &pkg_settings.SRequest{
		FReceiver: recv.ToString(),
		FHexData:  encoding.HexEncode(req.Bytes()),
	}
}
