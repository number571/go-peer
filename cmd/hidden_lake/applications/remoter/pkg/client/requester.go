package client

import (
	"context"
	"net/http"

	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	hls_request "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	"github.com/number571/go-peer/pkg/utils"
)

var (
	_ IRequester = &sRequester{}
)

type sRequester struct {
	fHLSClient hls_client.IClient
}

func NewRequester(pHLSClient hls_client.IClient) IRequester {
	return &sRequester{
		fHLSClient: pHLSClient,
	}
}

func (p *sRequester) Exec(pCtx context.Context, pAliasName string, pRequest hls_request.IRequest) ([]byte, error) {
	resp, err := p.fHLSClient.FetchRequest(pCtx, pAliasName, pRequest)
	if err != nil {
		return nil, utils.MergeErrors(ErrBadRequest, err)
	}

	if resp.GetCode() != http.StatusOK {
		return nil, ErrDecodeResponse
	}

	return resp.GetBody(), nil
}
