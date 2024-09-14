package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/distributor/pkg/config"
	hld_settings "github.com/number571/go-peer/cmd/hidden_lake/helpers/distributor/pkg/settings"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/response"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/utils"
)

const (
	cHandleIndexTemplate             = "%s" + hld_settings.CHandleIndexPath
	cHandleConfigSettingsTemplate    = "%s" + hld_settings.CHandleConfigSettings
	cHandleNetworkDistributeTemplate = "%s" + hld_settings.CHandleNetworkDistributePath
)

var (
	_ IRequester = &sRequester{}
)

type sRequester struct {
	fHost   string
	fClient *http.Client
}

func NewRequester(pHost string, pClient *http.Client) IRequester {
	return &sRequester{
		fHost:   pHost,
		fClient: pClient,
	}
}

func (p *sRequester) GetIndex(pCtx context.Context) (string, error) {
	res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleIndexTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return "", utils.MergeErrors(ErrBadRequest, err)
	}

	result := string(res)
	if result != hld_settings.CServiceFullName {
		return "", ErrInvalidTitle
	}

	return result, nil
}

func (p *sRequester) GetSettings(pCtx context.Context) (config.IConfigSettings, error) {
	res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleConfigSettingsTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return nil, utils.MergeErrors(ErrBadRequest, err)
	}

	cfgSettings := new(config.SConfigSettings)
	if err := encoding.DeserializeJSON(res, cfgSettings); err != nil {
		return nil, utils.MergeErrors(ErrDecodeResponse, err)
	}

	return cfgSettings, nil
}

func (p *sRequester) DistributeRequest(pCtx context.Context, pRequest request.IRequest) (response.IResponse, error) {
	res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleNetworkDistributeTemplate, p.fHost),
		pRequest.ToString(),
	)
	if err != nil {
		return nil, utils.MergeErrors(ErrBadRequest, err)
	}

	resp, err := response.LoadResponse(string(res))
	if err != nil {
		return nil, utils.MergeErrors(ErrDecodeResponse, err)
	}
	return resp, nil
}
