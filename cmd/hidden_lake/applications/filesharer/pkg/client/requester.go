package client

import (
	"context"
	"net/http"

	hlf_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/pkg/settings"
	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	hls_request "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/encoding"
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

func (p *sRequester) GetListFiles(pCtx context.Context, pAliasName string, pRequest hls_request.IRequest) ([]hlf_settings.SFileInfo, error) {
	resp, err := p.fHLSClient.FetchRequest(pCtx, pAliasName, pRequest)
	if err != nil {
		return nil, utils.MergeErrors(ErrRequest, err)
	}

	if resp.GetCode() != http.StatusOK {
		return nil, ErrDecodeResponse
	}

	list := make([]hlf_settings.SFileInfo, 0, hlf_settings.CDefaultPageOffset)
	if err := encoding.DeserializeJSON(resp.GetBody(), &list); err != nil {
		return nil, utils.MergeErrors(ErrInvalidResponse, err)
	}

	for _, info := range list {
		if len(encoding.HexDecode(info.FHash)) != hashing.CSHA256Size {
			return nil, ErrInvalidResponse
		}
	}

	return list, nil
}

func (p *sRequester) LoadFileChunk(pCtx context.Context, pAliasName string, pRequest hls_request.IRequest) ([]byte, error) {
	resp, err := p.fHLSClient.FetchRequest(pCtx, pAliasName, pRequest)
	if err != nil {
		return nil, utils.MergeErrors(ErrRequest, err)
	}

	if resp.GetCode() != http.StatusOK {
		return nil, ErrDecodeResponse
	}

	return resp.GetBody(), nil
}
