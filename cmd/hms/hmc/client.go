package hmc

import (
	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/errors"
	"github.com/number571/go-peer/local/payload"
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

func (client *sClient) Load(i uint64) (asymmetric.IPubKey, payload.IPayload, error) {
	msg, err := client.requester.Load(client.builder.Load(i))
	if err != nil {
		return nil, nil, errors.NewError(errors.CErrorServer, "message not loaded")
	}

	recv, pl := client.builder.(*sBuilder).client.Decrypt(msg)
	if recv == nil {
		return nil, nil, errors.NewError(errors.CErrorDecrypt, "message not decrypted")
	}

	return recv, pl, nil
}

func (client *sClient) Push(recv asymmetric.IPubKey, pl payload.IPayload) error {
	return client.requester.Push(client.builder.Push(recv, pl))
}
