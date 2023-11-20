package client

import (
	"fmt"

	net_message "github.com/number571/go-peer/pkg/network/message"
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
	res, err := p.fRequester.GetIndex()
	if err != nil {
		return "", fmt.Errorf("get index (client): %w", err)
	}
	return res, nil
}

func (p *sClient) GetHashes() ([]string, error) {
	res, err := p.fRequester.GetHashes()
	if err != nil {
		return nil, fmt.Errorf("get hashes (client): %w", err)
	}
	return res, nil
}

func (p *sClient) GetMessage(pHash string) (net_message.IMessage, error) {
	msg, err := p.fRequester.GetMessage(pHash)
	if err != nil {
		return nil, fmt.Errorf("get message (client): %w", err)
	}
	return msg, nil
}

func (p *sClient) PutMessage(pMsg net_message.IMessage) error {
	if err := p.fRequester.PutMessage(p.fBuilder.PutMessage(pMsg)); err != nil {
		return fmt.Errorf("put message (client): %w", err)
	}
	return nil
}
