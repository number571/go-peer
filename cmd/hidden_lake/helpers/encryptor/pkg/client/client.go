package client

import (
	"fmt"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	net_message "github.com/number571/go-peer/pkg/network/message"
)

var (
	_ IClient = &sClient{}
)

type sClient struct {
	fRequester IRequester
}

func NewClient(pRequester IRequester) IClient {
	return &sClient{
		fRequester: pRequester,
	}
}

func (p *sClient) GetIndex() (string, error) {
	res, err := p.fRequester.GetIndex()
	if err != nil {
		return "", fmt.Errorf("get index (client): %w", err)
	}
	return res, nil
}

func (p *sClient) EncryptMessage(pPubKey asymmetric.IPubKey, pData []byte) (net_message.IMessage, error) {
	res, err := p.fRequester.EncryptMessage(pPubKey, pData)
	if err != nil {
		return nil, fmt.Errorf("encrypt message (client): %w", err)
	}
	return res, nil
}

func (p *sClient) DecryptMessage(pNetMsg net_message.IMessage) (asymmetric.IPubKey, []byte, error) {
	pubKey, data, err := p.fRequester.DecryptMessage(pNetMsg)
	if err != nil {
		return nil, nil, fmt.Errorf("decrypt message (client): %w", err)
	}
	return pubKey, data, nil
}

func (p *sClient) GetPubKey() (asymmetric.IPubKey, error) {
	pubKey, err := p.fRequester.GetPubKey()
	if err != nil {
		return nil, fmt.Errorf("get public key (client): %w", err)
	}
	return pubKey, nil
}
