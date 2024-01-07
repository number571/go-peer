package client

import (
	"fmt"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/pkg/config"
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

func (p *sClient) GetPointer() (uint64, error) {
	res, err := p.fRequester.GetPointer()
	if err != nil {
		return 0, fmt.Errorf("get pointer (client): %w", err)
	}
	return res, nil
}

func (p *sClient) GetHash(i uint64) (string, error) {
	res, err := p.fRequester.GetHash(i)
	if err != nil {
		return "", fmt.Errorf("get hashes (client): %w", err)
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

func (p *sClient) GetSettings() (config.IConfigSettings, error) {
	res, err := p.fRequester.GetSettings()
	if err != nil {
		return nil, fmt.Errorf("get settings (client): %w", err)
	}
	return res, nil
}
