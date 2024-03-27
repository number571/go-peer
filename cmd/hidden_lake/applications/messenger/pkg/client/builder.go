package client

import (
	"net/http"

	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/pkg/settings"
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

func (p *sBuilder) PushMessage(pPseudonym string, pBody []byte) hls_request.IRequest {
	return hls_request.NewRequest(http.MethodPost, hlm_settings.CServiceFullName, hlm_settings.CPushPath).
		WithHead(map[string]string{
			hlm_settings.CHeaderPseudonym: pPseudonym,
		}).
		WithBody(pBody)
}
