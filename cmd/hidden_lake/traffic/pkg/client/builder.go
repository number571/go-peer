package client

import (
	"net/url"

	"github.com/number571/go-peer/pkg/client/message"
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

func (p *sBuilder) GetMessage(pHash string) string {
	return url.QueryEscape(pHash)
}

func (p *sBuilder) PutMessage(pMsg message.IMessage) string {
	return encoding.HexEncode(pMsg.ToBytes())
}
