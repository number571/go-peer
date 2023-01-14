package client

import (
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

var (
	_ IClient = &sClient{}
)

type sClient struct {
	fBuilder   IBuilder
	fRequester IRequester
}

func NewClient(requester IRequester) IClient {
	return &sClient{
		fBuilder:   NewBuilder(),
		fRequester: requester,
	}
}

func (client *sClient) Broadcast(recv asymmetric.IPubKey, data request.IRequest) error {
	return client.fRequester.Broadcast(client.fBuilder.Push(recv, data))
}

func (client *sClient) Request(recv asymmetric.IPubKey, data request.IRequest) ([]byte, error) {
	return client.fRequester.Request(client.fBuilder.Push(recv, data))
}

func (client *sClient) GetFriends() (map[string]asymmetric.IPubKey, error) {
	return client.fRequester.GetFriends()
}

func (client *sClient) AddFriend(aliasName string, pubKey asymmetric.IPubKey) error {
	return client.fRequester.AddFriend(client.fBuilder.Friend(aliasName, pubKey))
}

func (client *sClient) DelFriend(aliasName string) error {
	return client.fRequester.DelFriend(client.fBuilder.Friend(aliasName, nil))
}

func (client *sClient) GetOnlines() ([]string, error) {
	return client.fRequester.GetOnlines()
}

func (client *sClient) DelOnline(connect string) error {
	return client.fRequester.DelOnline(client.fBuilder.Connect(connect))
}

func (client *sClient) GetConnections() ([]string, error) {
	return client.fRequester.GetConnections()
}

func (client *sClient) AddConnection(connect string) error {
	return client.fRequester.AddConnection(client.fBuilder.Connect(connect))
}

func (client *sClient) DelConnection(connect string) error {
	return client.fRequester.DelConnection(client.fBuilder.Connect(connect))
}

func (client *sClient) PrivKey(privKey asymmetric.IPrivKey) error {
	return client.fRequester.PrivKey(client.fBuilder.PrivKey(privKey))
}

func (client *sClient) PubKey() (asymmetric.IPubKey, error) {
	return client.fRequester.PubKey()
}
