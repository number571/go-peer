package client

import (
	"context"
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

func (p *sClient) GetIndex(pCtx context.Context) (string, error) {
	res, err := p.fRequester.GetIndex(pCtx)
	if err != nil {
		return "", fmt.Errorf("get index (client): %w", err)
	}
	return res, nil
}

func (p *sClient) GetPointer(pCtx context.Context) (uint64, error) {
	res, err := p.fRequester.GetPointer(pCtx)
	if err != nil {
		return 0, fmt.Errorf("get pointer (client): %w", err)
	}
	return res, nil
}

func (p *sClient) GetHash(pCtx context.Context, i uint64) (string, error) {
	res, err := p.fRequester.GetHash(pCtx, i)
	if err != nil {
		return "", fmt.Errorf("get hashes (client): %w", err)
	}
	return res, nil
}

func (p *sClient) GetMessage(pCtx context.Context, pHash string) (net_message.IMessage, error) {
	msg, err := p.fRequester.GetMessage(pCtx, pHash)
	if err != nil {
		return nil, fmt.Errorf("get message (client): %w", err)
	}
	return msg, nil
}

func (p *sClient) PutMessage(pCtx context.Context, pMsg net_message.IMessage) error {
	if err := p.fRequester.PutMessage(pCtx, p.fBuilder.PutMessage(pMsg)); err != nil {
		return fmt.Errorf("put message (client): %w", err)
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
