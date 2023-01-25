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

func (client *sClient) GetIndex() (string, error) {
	return client.fRequester.GetIndex()
}

func (client *sClient) GetHashes() ([]string, error) {
	return client.fRequester.GetHashes()
}

func (client *sClient) GetMessage(hash string) (message.IMessage, error) {
	msg, err := client.fRequester.GetMessage(client.fBuilder.GetMessage(hash))
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (client *sClient) AddMessage(msg message.IMessage) error {
	return client.fRequester.AddMessage(client.fBuilder.AddMessage(msg))
}

func (client *sClient) DoBroadcast() error {
	return client.fRequester.DoBroadcast()
}
