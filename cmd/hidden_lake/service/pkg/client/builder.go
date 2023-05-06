package client

import (
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/client/message"
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

func (p *sBuilder) SetPrivKey(pPrivKey asymmetric.IPrivKey) *pkg_settings.SPrivKey {
	return &pkg_settings.SPrivKey{
		FPrivKey: pPrivKey.ToString(),
	}
}

func (p *sBuilder) Connect(pConnect string) *pkg_settings.SConnect {
	return &pkg_settings.SConnect{
		FConnect: pConnect,
	}
}

func (p *sBuilder) Friend(pAliasName string, pPubKey asymmetric.IPubKey) *pkg_settings.SFriend {
	if pPubKey == nil {
		return &pkg_settings.SFriend{
			FAliasName: pAliasName,
		}
	}
	return &pkg_settings.SFriend{
		FAliasName: pAliasName,
		FPublicKey: pPubKey.ToString(),
	}
}

func (p *sBuilder) Request(pRecv asymmetric.IPubKey, pReq request.IRequest) *pkg_settings.SRequest {
	return &pkg_settings.SRequest{
		FReceiver: pRecv.ToString(),
		FHexData:  encoding.HexEncode(pReq.ToBytes()),
	}
}

func (p *sBuilder) Message(pMsg message.IMessage) *pkg_settings.SMessage {
	return &pkg_settings.SMessage{
		FHexMessage: encoding.HexEncode(pMsg.ToBytes()),
	}
}
