package client

import (
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/response"
	"github.com/number571/go-peer/pkg/client/message"
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
	return p.fRequester.GetIndex()
}

func (p *sClient) HandleMessage(pMsg message.IMessage) error {
	return p.fRequester.HandleMessage(p.fBuilder.Message(pMsg))
}

func (p *sClient) BroadcastRequest(pRecv asymmetric.IPubKey, pData request.IRequest) error {
	return p.fRequester.BroadcastRequest(p.fBuilder.Request(pRecv, pData))
}

func (p *sClient) FetchRequest(pRecv asymmetric.IPubKey, pData request.IRequest) (response.IResponse, error) {
	return p.fRequester.FetchRequest(p.fBuilder.Request(pRecv, pData))
}

func (p *sClient) GetFriends() (map[string]asymmetric.IPubKey, error) {
	return p.fRequester.GetFriends()
}

func (p *sClient) AddFriend(pAliasName string, pPubKey asymmetric.IPubKey) error {
	return p.fRequester.AddFriend(p.fBuilder.Friend(pAliasName, pPubKey))
}

func (p *sClient) DelFriend(pAliasName string) error {
	return p.fRequester.DelFriend(p.fBuilder.Friend(pAliasName, nil))
}

func (p *sClient) GetOnlines() ([]string, error) {
	return p.fRequester.GetOnlines()
}

func (p *sClient) DelOnline(pConnect string) error {
	return p.fRequester.DelOnline(p.fBuilder.Connect(pConnect))
}

func (p *sClient) GetConnections() ([]string, error) {
	return p.fRequester.GetConnections()
}

func (p *sClient) AddConnection(pConnect string) error {
	return p.fRequester.AddConnection(p.fBuilder.Connect(pConnect))
}

func (p *sClient) DelConnection(pConnect string) error {
	return p.fRequester.DelConnection(p.fBuilder.Connect(pConnect))
}

func (p *sClient) SetPrivKey(pPrivKey asymmetric.IPrivKey) error {
	return p.fRequester.SetPrivKey(p.fBuilder.SetPrivKey(pPrivKey))
}

func (p *sClient) GetPubKey() (asymmetric.IPubKey, error) {
	return p.fRequester.GetPubKey()
}
