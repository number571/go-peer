package client

import "context"

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

func (p *sClient) PushMessage(pCtx context.Context, pAliasName string, pPseudonym string, pBody []byte) error {
	return p.fRequester.PushMessage(pCtx, pAliasName, p.fBuilder.PushMessage(pPseudonym, pBody))
}
