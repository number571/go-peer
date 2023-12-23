package client

import (
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

func (p *sClient) SetNetworkKey(pNetworkKey string) error {
	err := p.fRequester.SetNetworkKey(pNetworkKey)
	if err != nil {
		return fmt.Errorf("set network key (client): %w", err)
	}
	return nil
}

func (p *sClient) BroadcastRequest(pRecv string, pData request.IRequest) error {
	if err := p.fRequester.BroadcastRequest(p.fBuilder.Request(pRecv, pData)); err != nil {
		return fmt.Errorf("broadcast request (client): %w", err)
	}
	return nil
}

func (p *sClient) FetchRequest(pRecv string, pData request.IRequest) (response.IResponse, error) {
	res, err := p.fRequester.FetchRequest(p.fBuilder.Request(pRecv, pData))
	if err != nil {
		return nil, fmt.Errorf("fetch request (client): %w", err)
	}
	return res, nil
}

func (p *sClient) GetFriends() (map[string]asymmetric.IPubKey, error) {
	res, err := p.fRequester.GetFriends()
	if err != nil {
		return nil, fmt.Errorf("get friends (client): %w", err)
	}
	return res, nil
}

func (p *sClient) AddFriend(pAliasName string, pPubKey asymmetric.IPubKey) error {
	if err := p.fRequester.AddFriend(p.fBuilder.Friend(pAliasName, pPubKey)); err != nil {
		return fmt.Errorf("add friend (client): %w", err)
	}
	return nil
}

func (p *sClient) DelFriend(pAliasName string) error {
	if err := p.fRequester.DelFriend(p.fBuilder.Friend(pAliasName, nil)); err != nil {
		return fmt.Errorf("del friend (client): %w", err)
	}
	return nil
}

func (p *sClient) GetOnlines() ([]string, error) {
	res, err := p.fRequester.GetOnlines()
	if err != nil {
		return nil, fmt.Errorf("get onlines (client): %w", err)
	}
	return res, nil
}

func (p *sClient) DelOnline(pConnect string) error {
	if err := p.fRequester.DelOnline(pConnect); err != nil {
		return fmt.Errorf("del online (client): %w", err)
	}
	return nil
}

func (p *sClient) GetConnections() ([]string, error) {
	res, err := p.fRequester.GetConnections()
	if err != nil {
		return nil, fmt.Errorf("get connections (client): %w", err)
	}
	return res, nil
}

func (p *sClient) AddConnection(pConnect string) error {
	if err := p.fRequester.AddConnection(pConnect); err != nil {
		return fmt.Errorf("add connection (client): %w", err)
	}
	return nil
}

func (p *sClient) DelConnection(pConnect string) error {
	if err := p.fRequester.DelConnection(pConnect); err != nil {
		return fmt.Errorf("del connection (client): %w", err)
	}
	return nil
}

func (p *sClient) GetPubKey() (asymmetric.IPubKey, error) {
	pubKey, err := p.fRequester.GetPubKey()
	if err != nil {
		return nil, fmt.Errorf("get public key (client): %w", err)
	}
	return pubKey, nil
}
