package hmc

import (
	"fmt"

	hms_settings "github.com/number571/go-peer/cmd/hms/settings"
	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/local"
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

func (client *sClient) Load(i uint64) (local.IMessage, error) {
	msg, err := client.requester.Load(client.builder.Load(i))
	if err != nil {
		return nil, err
	}

	msg, title := client.builder.(*sBuiler).client.Decrypt(msg)
	if msg == nil {
		return nil, fmt.Errorf("message is nil")
	}

	if string(title) != hms_settings.CTitlePattern {
		return nil, fmt.Errorf("title is not equal")
	}

	return msg, nil
}

func (client *sClient) Push(receiver crypto.IPubKey, body []byte) error {
	return client.requester.Push(client.builder.Push(receiver, body))
}
