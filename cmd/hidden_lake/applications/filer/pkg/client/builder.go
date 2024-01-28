package client

import (
	"fmt"
	"net/http"

	hlf_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/filer/pkg/settings"
	hls_request "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
)

var (
	_ IBuilder = &sBuilder{}
)

type sBuilder struct {
}

func NewBuilder() IBuilder {
	return &sBuilder{}
}

func (p *sBuilder) GetListFiles(pPage uint64) hls_request.IRequest {
	return hls_request.NewRequest(
		http.MethodGet,
		hlf_settings.CTitlePattern,
		fmt.Sprintf("%s?page=%d", hlf_settings.CListPath, pPage),
	)
}

func (p *sBuilder) LoadFileChunk(pName string, pChunk uint64) hls_request.IRequest {
	return hls_request.NewRequest(
		http.MethodGet,
		hlf_settings.CTitlePattern,
		fmt.Sprintf("%s?name=%s&chunk=%d", hlf_settings.CLoadPath, pName, pChunk),
	)
}
