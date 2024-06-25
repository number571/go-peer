package client

import (
	"context"

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

func (p *sRequester) PushMessage(pCtx context.Context, pAliasName string, pRequest hls_request.IRequest) error {
	if err := p.fHLSClient.BroadcastRequest(pCtx, pAliasName, pRequest); err != nil {
		return utils.MergeErrors(ErrPushMessage, err)
	}
	return nil
}
