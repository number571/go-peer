package client

import (
	"fmt"
)

var (
	_ IClient = &sClient{}
)

type sClient struct {
	fRequester IRequester
}

func NewClient(pRequester IRequester) IClient {
	return &sClient{
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

func (p *sClient) RunTransfer() error {
	if err := p.fRequester.RunTransfer(); err != nil {
		return fmt.Errorf("run loader (client): %w", err)
	}
	return nil
}

func (p *sClient) StopTransfer() error {
	if err := p.fRequester.StopTransfer(); err != nil {
		return fmt.Errorf("stop loader (client): %w", err)
	}
	return nil
}
