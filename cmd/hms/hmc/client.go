package hmc

import (
	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/message"
	"github.com/number571/go-peer/modules/payload"
)

var (
	_ IClient = &sClient{}
)

type sClient struct {
	builder   IBuilder
	requester IRequester
}

func NewClient(builder IBuilder, requester IRequester) IClient {
	return &sClient{
		builder:   builder,
		requester: requester,
	}
}

func (client *sClient) Size() (uint64, error) {
	return client.requester.Size(client.builder.Size())
}

func (client *sClient) Load(i uint64) (message.IMessage, error) {
	msg, err := client.requester.Load(client.builder.Load(i))
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (client *sClient) Push(recv asymmetric.IPubKey, pl payload.IPayload) error {
	return client.requester.Push(client.builder.Push(recv, pl))
}
