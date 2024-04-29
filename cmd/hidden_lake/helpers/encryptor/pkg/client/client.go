package client

import (
	"context"
	"fmt"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/encryptor/pkg/config"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
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

func (p *sClient) EncryptMessage(pCtx context.Context, pPubKey asymmetric.IPubKey, pPayload payload.IPayload) (net_message.IMessage, error) {
	res, err := p.fRequester.EncryptMessage(pCtx, pPubKey, pPayload)
	if err != nil {
		return nil, fmt.Errorf("encrypt message (client): %w", err)
	}
	return res, nil
}

func (p *sClient) DecryptMessage(pCtx context.Context, pNetMsg net_message.IMessage) (asymmetric.IPubKey, payload.IPayload, error) {
	pubKey, data, err := p.fRequester.DecryptMessage(pCtx, pNetMsg)
	if err != nil {
		return nil, nil, fmt.Errorf("decrypt message (client): %w", err)
	}
	return pubKey, data, nil
}

func (p *sClient) GetPubKey(pCtx context.Context) (asymmetric.IPubKey, error) {
	pubKey, err := p.fRequester.GetPubKey(pCtx)
	if err != nil {
		return nil, fmt.Errorf("get public key (client): %w", err)
	}
	return pubKey, nil
}

func (p *sClient) GetSettings(pCtx context.Context) (config.IConfigSettings, error) {
	res, err := p.fRequester.GetSettings(pCtx)
	if err != nil {
		return nil, fmt.Errorf("get settings (client): %w", err)
	}
	return res, nil
}
