package client

import (
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
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

func (p *sBuilder) Friend(pAliasName string, pSharedKey []byte) *pkg_settings.SFriend {
	if pSharedKey == nil {
		// del friend
		return &pkg_settings.SFriend{
			FAliasName: pAliasName,
		}
	}
	// add friend
	return &pkg_settings.SFriend{
		FAliasName: pAliasName,
		FSharedKey: encoding.HexEncode(pSharedKey),
	}
}

func (p *sBuilder) Request(pReceiver string, pReq request.IRequest) *pkg_settings.SRequest {
	return &pkg_settings.SRequest{
		FReceiver: pReceiver,
		FReqData:  pReq.(*request.SRequest),
	}
}
