package client

import (
	"context"

	hlf_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/pkg/settings"
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

func (p *sClient) GetListFiles(pCtx context.Context, pAliasName string, pPage uint64) ([]hlf_settings.SFileInfo, error) {
	return p.fRequester.GetListFiles(pCtx, pAliasName, p.fBuilder.GetListFiles(pPage))
}

func (p *sClient) LoadFileChunk(pCtx context.Context, pAliasName, pName string, pChunk uint64) ([]byte, error) {
	return p.fRequester.LoadFileChunk(pCtx, pAliasName, p.fBuilder.LoadFileChunk(pName, pChunk))
}
