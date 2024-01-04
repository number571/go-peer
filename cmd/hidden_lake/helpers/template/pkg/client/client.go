package client

import (
	"fmt"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/template/pkg/config"
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

func (p *sClient) GetSettings() (config.IConfigSettings, error) {
	res, err := p.fRequester.GetSettings()
	if err != nil {
		return nil, fmt.Errorf("get settings (client): %w", err)
	}
	return res, nil
}
