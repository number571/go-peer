package client

import (
	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	hls_request "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
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

func (p *sRequester) PushMessage(pAliasName string, pRequest hls_request.IRequest) error {
	if err := p.fHLSClient.BroadcastRequest(pAliasName, pRequest); err != nil {
		return err // TODO: create errors
	}
	return nil
}
