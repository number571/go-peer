package client

import (
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/encoding"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
)

var (
	_ IBuilder = &sBuilder{}
)

type sBuilder struct {
}

func NewBuilder() IBuilder {
	return &sBuilder{}
}

func (p *sBuilder) GetMessage(pHash string) *pkg_settings.SLoadRequest {
	return &pkg_settings.SLoadRequest{
		FHash: pHash,
	}
}

func (p *sBuilder) PutMessage(pMsg message.IMessage) *pkg_settings.SPushRequest {
	return &pkg_settings.SPushRequest{
		FMessage: encoding.HexEncode(pMsg.ToBytes()),
	}
}
