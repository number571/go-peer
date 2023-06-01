package client

import (
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/errors"
)

var (
	_ IClient = &sClient{}
)

type sClient struct {
	fBuilder   IBuilder
	fRequester IRequester
}

func NewClient(pBuilder IBuilder, pRequester IRequester) IClient {
	return &sClient{
		fBuilder:   pBuilder,
		fRequester: pRequester,
	}
}

func (p *sClient) GetIndex() (string, error) {
	res, err := p.fRequester.GetIndex()
	if err != nil {
		return "", errors.WrapError(err, "get index (client)")
	}
	return res, nil
}

func (p *sClient) GetHashes() ([]string, error) {
	res, err := p.fRequester.GetHashes()
	if err != nil {
		return nil, errors.WrapError(err, "get hashes (client)")
	}
	return res, nil
}

func (p *sClient) GetMessage(pHash string) (message.IMessage, error) {
	msg, err := p.fRequester.GetMessage(p.fBuilder.GetMessage(pHash))
	if err != nil {
		return nil, errors.WrapError(err, "get message (client)")
	}
	return msg, nil
}

func (p *sClient) PutMessage(pMsg message.IMessage) error {
	if err := p.fRequester.PutMessage(p.fBuilder.PutMessage(pMsg)); err != nil {
		return errors.WrapError(err, "put message (client)")
	}
	return nil
}
