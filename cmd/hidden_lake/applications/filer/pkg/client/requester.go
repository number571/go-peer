package client

import (
	"errors"
	"fmt"
	"net/http"

	hlf_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/filer/pkg/settings"
	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	hls_request "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/encoding"
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

func (p *sRequester) GetListFiles(pAliasName string, pRequest hls_request.IRequest) ([]hlf_settings.SFileInfo, error) {
	resp, err := p.fHLSClient.FetchRequest(pAliasName, pRequest)
	if err != nil {
		return nil, err // TODO: create errors
	}
	if resp.GetCode() != http.StatusOK {
		return nil, fmt.Errorf("got %d code", resp.GetCode()) // TODO: create errors
	}
	list := make([]hlf_settings.SFileInfo, 0, hlf_settings.CPageOffset)
	if err := encoding.DeserializeJSON(resp.GetBody(), &list); err != nil {
		return nil, err // TODO: create errors
	}
	for _, info := range list {
		if len(encoding.HexDecode(info.FHash)) != hashing.CSHA256Size {
			return nil, errors.New("got invalid hash value")
		}
	}
	return list, nil
}

func (p *sRequester) LoadFileChunk(pAliasName string, pRequest hls_request.IRequest) ([]byte, error) {
	resp, err := p.fHLSClient.FetchRequest(pAliasName, pRequest)
	if err != nil {
		return nil, err // TODO: create errors
	}
	if resp.GetCode() != http.StatusOK {
		return nil, fmt.Errorf("got %d code", resp.GetCode()) // TODO: create errors
	}
	return resp.GetBody(), nil
}
