package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/pkg/config"
	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/encoding"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/utils"
)

var (
	_ IRequester = &sRequester{}
)

const (
	cHandleIndexTemplate          = "%s" + hlt_settings.CHandleIndexPath
	cHandleStoragePointerTemplate = "%s" + hlt_settings.CHandleStoragePointerPath
	cHandleStorageHashesTemplate  = "%s" + hlt_settings.CHandleStorageHashesPath + "?id=%d"
	cHandleNetworkMessageTemplate = "%s" + hlt_settings.CHandleNetworkMessagePath
	cHandleConfigSettingsTemplate = "%s" + hlt_settings.CHandleConfigSettings
)

type sRequester struct {
	fHost   string
	fClient *http.Client
	fParams net_message.ISettings
}

func NewRequester(pHost string, pClient *http.Client, pParams net_message.ISettings) IRequester {
	return &sRequester{
		fHost:   pHost,
		fClient: pClient,
		fParams: pParams,
	}
}

func (p *sRequester) GetIndex(pCtx context.Context) (string, error) {
	resp, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleIndexTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return "", utils.MergeErrors(ErrBadRequest, err)
	}

	result := string(resp)
	if result != hlt_settings.CServiceFullName {
		return "", ErrInvalidTitle
	}

	return result, nil
}

func (p *sRequester) GetPointer(pCtx context.Context) (uint64, error) {
	resp, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleStoragePointerTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return 0, utils.MergeErrors(ErrBadRequest, err)
	}

	pointer, err := strconv.ParseUint(string(resp), 10, 64)
	if err != nil {
		return 0, utils.MergeErrors(ErrDecodeResponse, err)
	}

	return uint64(pointer), nil
}

func (p *sRequester) GetHash(pCtx context.Context, i uint64) (string, error) {
	resp, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleStorageHashesTemplate, p.fHost, i),
		nil,
	)
	if err != nil {
		return "", utils.MergeErrors(ErrBadRequest, err)
	}

	// response in hex encoding
	if len(resp) != 2*hashing.CSHA256Size {
		return "", ErrDecodeResponse
	}

	return string(resp), nil
}

func (p *sRequester) GetMessage(pCtx context.Context, pHash string) (net_message.IMessage, error) {
	resp, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleNetworkMessageTemplate+"?hash=%s", p.fHost, pHash),
		nil,
	)
	if err != nil {
		return nil, utils.MergeErrors(ErrBadRequest, err)
	}

	msg, err := net_message.LoadMessage(p.fParams, string(resp))
	if err != nil {
		return nil, utils.MergeErrors(ErrDecodeMessage, err)
	}

	if !bytes.Equal(msg.GetHash(), encoding.HexDecode(pHash)) {
		return nil, ErrInvalidHexFormat
	}

	return msg, nil
}

func (p *sRequester) PutMessage(pCtx context.Context, pRequest string) error {
	_, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleNetworkMessageTemplate, p.fHost),
		pRequest,
	)
	if err != nil {
		return utils.MergeErrors(ErrBadRequest, err)
	}
	return nil
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
