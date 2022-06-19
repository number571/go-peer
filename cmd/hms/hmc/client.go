package hmc

import (
	hms_settings "github.com/number571/go-peer/cmd/hms/settings"
	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/errors"
	"github.com/number571/go-peer/local/message"
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
		return nil, errors.NewError(errors.CErrorServer, "message not loaded")
	}

	msg, title := client.builder.(*sBuiler).client.Decrypt(msg)
	if msg == nil {
		return nil, errors.NewError(errors.CErrorDecrypt, "message not decrypted")
	}

	if string(title) != hms_settings.CTitlePattern {
		return nil, errors.NewError(errors.CErrorNotEqual, "title is not equal")
	}

	return msg, nil
}

func (client *sClient) Push(receiver asymmetric.IPubKey, body []byte) error {
	return client.requester.Push(client.builder.Push(receiver, body))
}
