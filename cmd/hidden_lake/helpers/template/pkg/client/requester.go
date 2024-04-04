package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/template/pkg/config"
	hl_t_settings "github.com/number571/go-peer/cmd/hidden_lake/helpers/template/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/utils"
)

const (
	cHandleIndexTemplate          = "%s" + hl_t_settings.CHandleIndexPath
	cHandleConfigSettingsTemplate = "%s" + hl_t_settings.CHandleConfigSettingsPath
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
	if result != hl_t_settings.CServiceFullName {
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
	if err := encoding.DeserializeJSON([]byte(res), cfgSettings); err != nil {
		return nil, utils.MergeErrors(ErrDecodeResponse, err)
	}

	return cfgSettings, nil
}
