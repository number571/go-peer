package client

import (
	"context"
	"fmt"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/loader/pkg/config"
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

func (p *sClient) GetIndex(pCtx context.Context) (string, error) {
	res, err := p.fRequester.GetIndex(pCtx)
	if err != nil {
		return "", fmt.Errorf("get index (client): %w", err)
	}
	return res, nil
}

func (p *sClient) RunTransfer(pCtx context.Context) error {
	if err := p.fRequester.RunTransfer(pCtx); err != nil {
		return fmt.Errorf("run loader (client): %w", err)
	}
	return nil
}

func (p *sClient) StopTransfer(pCtx context.Context) error {
	if err := p.fRequester.StopTransfer(pCtx); err != nil {
		return fmt.Errorf("stop loader (client): %w", err)
	}
	return nil
}

func (p *sClient) GetSettings(pCtx context.Context) (config.IConfigSettings, error) {
	res, err := p.fRequester.GetSettings(pCtx)
	if err != nil {
		return nil, fmt.Errorf("get settings (client): %w", err)
	}
	return res, nil
}
