package hlc

import (
	hls_network "github.com/number571/go-peer/cmd/hls/network"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
)

var (
	_ IClient = &sClient{}
)

type sClient struct {
	builder   IBuilder
	requester IRequester
}

func NewClient(requester IRequester) IClient {
	return &sClient{
		builder:   NewBuilder(),
		requester: requester,
	}
}

func (client *sClient) Broadcast(recv asymmetric.IPubKey, data hls_network.IRequest) error {
	return client.requester.Broadcast(client.builder.Push(recv, data))
}

func (client *sClient) Request(recv asymmetric.IPubKey, data hls_network.IRequest) ([]byte, error) {
	return client.requester.Request(client.builder.Push(recv, data))
}

func (client *sClient) GetFriends() (map[string]asymmetric.IPubKey, error) {
	return client.requester.GetFriends()
}

func (client *sClient) AddFriend(aliasName string, pubKey asymmetric.IPubKey) error {
	return client.requester.AddFriend(client.builder.Friend(aliasName, pubKey))
}

func (client *sClient) DelFriend(aliasName string) error {
	return client.requester.DelFriend(client.builder.Friend(aliasName, nil))
}

func (client *sClient) GetOnlines() ([]string, error) {
	return client.requester.GetOnlines()
}

func (client *sClient) DelOnline(connect string) error {
	return client.requester.DelOnline(client.builder.Connect(connect))
}

func (client *sClient) GetConnections() ([]string, error) {
	return client.requester.GetConnections()
}

func (client *sClient) AddConnection(connect string) error {
	return client.requester.AddConnection(client.builder.Connect(connect))
}

func (client *sClient) DelConnection(connect string) error {
	return client.requester.DelConnection(client.builder.Connect(connect))
}

func (client *sClient) PubKey() (asymmetric.IPubKey, error) {
	return client.requester.PubKey()
}
