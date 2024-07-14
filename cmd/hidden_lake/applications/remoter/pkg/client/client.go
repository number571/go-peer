package client

import (
	"context"
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

func (p *sClient) Exec(pCtx context.Context, pAliasName string, pCmd ...string) ([]byte, error) {
	return p.fRequester.Exec(pCtx, pAliasName, p.fBuilder.Exec(pCmd...))
}
