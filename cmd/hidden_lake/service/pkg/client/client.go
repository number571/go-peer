package client

import (
	"context"
	"fmt"

	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/config"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/response"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
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

func (p *sClient) GetSettings(pCtx context.Context) (config.IConfigSettings, error) {
	res, err := p.fRequester.GetSettings(pCtx)
	if err != nil {
		return nil, fmt.Errorf("get settings (client): %w", err)
	}
	return res, nil
}

func (p *sClient) SetNetworkKey(pCtx context.Context, pNetworkKey string) error {
	err := p.fRequester.SetNetworkKey(pCtx, pNetworkKey)
	if err != nil {
		return fmt.Errorf("set network key (client): %w", err)
	}
	return nil
}

func (p *sClient) BroadcastRequest(pCtx context.Context, pRecv string, pData request.IRequest) error {
	if err := p.fRequester.BroadcastRequest(pCtx, p.fBuilder.Request(pRecv, pData)); err != nil {
		return fmt.Errorf("broadcast request (client): %w", err)
	}
	return nil
}

func (p *sClient) FetchRequest(pCtx context.Context, pRecv string, pData request.IRequest) (response.IResponse, error) {
	res, err := p.fRequester.FetchRequest(pCtx, p.fBuilder.Request(pRecv, pData))
	if err != nil {
		return nil, fmt.Errorf("fetch request (client): %w", err)
	}
	return res, nil
}

func (p *sClient) GetFriends(pCtx context.Context) (map[string]asymmetric.IPubKey, error) {
	res, err := p.fRequester.GetFriends(pCtx)
	if err != nil {
		return nil, fmt.Errorf("get friends (client): %w", err)
	}
	return res, nil
}

func (p *sClient) AddFriend(pCtx context.Context, pAliasName string, pPubKey asymmetric.IPubKey) error {
	if err := p.fRequester.AddFriend(pCtx, p.fBuilder.Friend(pAliasName, pPubKey)); err != nil {
		return fmt.Errorf("add friend (client): %w", err)
	}
	return nil
}

func (p *sClient) DelFriend(pCtx context.Context, pAliasName string) error {
	if err := p.fRequester.DelFriend(pCtx, p.fBuilder.Friend(pAliasName, nil)); err != nil {
		return fmt.Errorf("del friend (client): %w", err)
	}
	return nil
}

func (p *sClient) GetOnlines(pCtx context.Context) ([]string, error) {
	res, err := p.fRequester.GetOnlines(pCtx)
	if err != nil {
		return nil, fmt.Errorf("get onlines (client): %w", err)
	}
	return res, nil
}

func (p *sClient) DelOnline(pCtx context.Context, pConnect string) error {
	if err := p.fRequester.DelOnline(pCtx, pConnect); err != nil {
		return fmt.Errorf("del online (client): %w", err)
	}
	return nil
}

func (p *sClient) GetConnections(pCtx context.Context) ([]string, error) {
	res, err := p.fRequester.GetConnections(pCtx)
	if err != nil {
		return nil, fmt.Errorf("get connections (client): %w", err)
	}
	return res, nil
}

func (p *sClient) AddConnection(pCtx context.Context, pConnect string) error {
	if err := p.fRequester.AddConnection(pCtx, pConnect); err != nil {
		return fmt.Errorf("add connection (client): %w", err)
	}
	return nil
}

func (p *sClient) DelConnection(pCtx context.Context, pConnect string) error {
	if err := p.fRequester.DelConnection(pCtx, pConnect); err != nil {
		return fmt.Errorf("del connection (client): %w", err)
	}
	return nil
}

func (p *sClient) GetPubKey(pCtx context.Context) (asymmetric.IPubKey, error) {
	pubKey, err := p.fRequester.GetPubKey(pCtx)
	if err != nil {
		return nil, fmt.Errorf("get public key (client): %w", err)
	}
	return pubKey, nil
}
