package client

import (
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/response"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/errors"
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
		return "", errors.WrapError(err, "get index (client)")
	}
	return res, nil
}

func (p *sClient) GetNetworkKey() (string, error) {
	res, err := p.fRequester.GetNetworkKey()
	if err != nil {
		return "", errors.WrapError(err, "get network key (client)")
	}
	return res, nil
}

func (p *sClient) SetNetworkKey(pNetworkKey string) error {
	err := p.fRequester.SetNetworkKey(pNetworkKey)
	if err != nil {
		return errors.WrapError(err, "set network key (client)")
	}
	return nil
}

func (p *sClient) HandleMessage(pMsg message.IMessage) error {
	if err := p.fRequester.HandleMessage(p.fBuilder.Message(pMsg)); err != nil {
		return errors.WrapError(err, "handle message (client)")
	}
	return nil
}

func (p *sClient) BroadcastRequest(pRecv string, pData request.IRequest) error {
	if err := p.fRequester.BroadcastRequest(p.fBuilder.Request(pRecv, pData)); err != nil {
		return errors.WrapError(err, "broadcast request (client)")
	}
	return nil
}

func (p *sClient) FetchRequest(pRecv string, pData request.IRequest) (response.IResponse, error) {
	res, err := p.fRequester.FetchRequest(p.fBuilder.Request(pRecv, pData))
	if err != nil {
		return nil, errors.WrapError(err, "fetch request (client)")
	}
	return res, nil
}

func (p *sClient) GetFriends() (map[string]asymmetric.IPubKey, error) {
	res, err := p.fRequester.GetFriends()
	if err != nil {
		return nil, errors.WrapError(err, "get friends (client)")
	}
	return res, nil
}

func (p *sClient) AddFriend(pAliasName string, pPubKey asymmetric.IPubKey) error {
	if err := p.fRequester.AddFriend(p.fBuilder.Friend(pAliasName, pPubKey)); err != nil {
		return errors.WrapError(err, "add friend (client)")
	}
	return nil
}

func (p *sClient) DelFriend(pAliasName string) error {
	if err := p.fRequester.DelFriend(p.fBuilder.Friend(pAliasName, nil)); err != nil {
		return errors.WrapError(err, "del friend (client)")
	}
	return nil
}

func (p *sClient) GetOnlines() ([]string, error) {
	res, err := p.fRequester.GetOnlines()
	if err != nil {
		return nil, errors.WrapError(err, "get onlines (client)")
	}
	return res, nil
}

func (p *sClient) DelOnline(pConnect string) error {
	if err := p.fRequester.DelOnline(pConnect); err != nil {
		return errors.WrapError(err, "del online (client)")
	}
	return nil
}

func (p *sClient) GetConnections() ([]string, error) {
	res, err := p.fRequester.GetConnections()
	if err != nil {
		return nil, errors.WrapError(err, "get connections (client)")
	}
	return res, nil
}

func (p *sClient) AddConnection(pConnect string) error {
	if err := p.fRequester.AddConnection(pConnect); err != nil {
		return errors.WrapError(err, "add connection (client)")
	}
	return nil
}

func (p *sClient) DelConnection(pConnect string) error {
	if err := p.fRequester.DelConnection(pConnect); err != nil {
		return errors.WrapError(err, "del connection (client)")
	}
	return nil
}

func (p *sClient) SetPrivKey(pEphPubKey asymmetric.IPubKey, pPrivKey asymmetric.IPrivKey) error {
	if err := p.fRequester.SetPrivKey(p.fBuilder.SetPrivKey(pEphPubKey, pPrivKey)); err != nil {
		return errors.WrapError(err, "set private key (client)")
	}
	return nil
}

func (p *sClient) GetPubKey() (asymmetric.IPubKey, asymmetric.IPubKey, error) {
	pubKey, ephPubKey, err := p.fRequester.GetPubKey()
	if err != nil {
		return nil, nil, errors.WrapError(err, "get public key (client)")
	}
	return pubKey, ephPubKey, nil
}
