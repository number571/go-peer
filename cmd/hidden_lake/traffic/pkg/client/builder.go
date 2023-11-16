package client

import net_message "github.com/number571/go-peer/pkg/network/message"

var (
	_ IBuilder = &sBuilder{}
)

type sBuilder struct {
}

func NewBuilder() IBuilder {
	return &sBuilder{}
}

func (p *sBuilder) PutMessage(pMsg net_message.IMessage) string {
	return pMsg.ToString()
}
