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

func (client *sClient) Request(recv asymmetric.IPubKey, data hls_network.IRequest) ([]byte, error) {
	return client.requester.Request(client.builder.Request(recv, data))
}

func (client *sClient) Friends() ([]asymmetric.IPubKey, error) {
	return client.requester.Friends()
}

func (client *sClient) Online() ([]string, error) {
	return client.requester.Online()
}

func (client *sClient) PubKey() (asymmetric.IPubKey, error) {
	return client.requester.PubKey()
}
