package client

import (
	"github.com/number571/go-peer/pkg/client/message"
)

var (
	_ IClient = &sClient{}
)

type sClient struct {
	fBuilder   IBuilder
	fRequester IRequester
}

func NewClient(pBuilder IBuilder, pRequester IRequester) IClient {
	return &sClient{
		fBuilder:   pBuilder,
		fRequester: pRequester,
	}
}

func (p *sClient) GetIndex() (string, error) {
	return p.fRequester.GetIndex()
}

func (p *sClient) GetHashes() ([]string, error) {
	return p.fRequester.GetHashes()
}

func (p *sClient) GetMessage(pHash string) (message.IMessage, error) {
	msg, err := p.fRequester.GetMessage(p.fBuilder.GetMessage(pHash))
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (p *sClient) PutMessage(pMsg message.IMessage) error {
	return p.fRequester.PutMessage(p.fBuilder.PutMessage(pMsg))
}
