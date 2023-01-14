package client

import (
	"github.com/number571/go-peer/pkg/client/message"
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

func (client *sClient) Hashes() ([]string, error) {
	return client.fRequester.Hashes()
}

func (client *sClient) Load(hash string) (message.IMessage, error) {
	msg, err := client.fRequester.Load(client.fBuilder.Load(hash))
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (client *sClient) Push(msg message.IMessage) error {
	return client.fRequester.Push(client.fBuilder.Push(msg))
}
