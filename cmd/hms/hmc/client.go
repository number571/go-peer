package hmc

import (
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/payload"
)

var (
	_ IClient = &sClient{}
)

type sClient struct {
	fBuilder   IBuilder
	fRequester IRequester
}

func NewClient(builder IBuilder, requester IRequester) IClient {
	return &sClient{
		fBuilder:   builder,
		fRequester: requester,
	}
}

func (client *sClient) Size() (uint64, error) {
	return client.fRequester.Size(client.fBuilder.Size())
}

func (client *sClient) Load(i uint64) (message.IMessage, error) {
	msg, err := client.fRequester.Load(client.fBuilder.Load(i))
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (client *sClient) Push(recv asymmetric.IPubKey, pl payload.IPayload) error {
	return client.fRequester.Push(client.fBuilder.Push(recv, pl))
}
