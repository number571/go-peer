package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/encryptor/pkg/config"
	hle_settings "github.com/number571/go-peer/cmd/hidden_lake/helpers/encryptor/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/utils"
)

const (
	cHandleIndexTemplate          = "%s" + hle_settings.CHandleIndexPath
	cHandleMessageEncryptTemplate = "%s" + hle_settings.CHandleMessageEncryptPath
	cHandleMessageDecryptTemplate = "%s" + hle_settings.CHandleMessageDecryptPath
	cHandleServicePubKeyTemplate  = "%s" + hle_settings.CHandleServicePubKeyPath
	cHandleConfigSettingsTemplate = "%s" + hle_settings.CHandleConfigSettings
)

var (
	_ IRequester = &sRequester{}
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
	res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleIndexTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return "", utils.MergeErrors(ErrRequest, err)
	}

	result := string(res)
	if result != hle_settings.CServiceFullName {
		return "", utils.MergeErrors(ErrDecodeResponse, errors.New("incorrect title pattern"))
	}

	return result, nil
}

func (p *sRequester) EncryptMessage(pCtx context.Context, pPubKey asymmetric.IPubKey, pData []byte) (net_message.IMessage, error) {
	resp, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleMessageEncryptTemplate, p.fHost),
		hle_settings.SContainer{
			FPublicKey: pPubKey.ToString(),
			FHexData:   encoding.HexEncode(pData),
		},
	)
	if err != nil {
		return nil, utils.MergeErrors(ErrRequest, err)
	}

	msg, err := net_message.LoadMessage(p.fParams, string(resp))
	if err != nil {
		return nil, utils.MergeErrors(ErrDecodeResponse, err)
	}

	return msg, nil
}

func (p *sRequester) DecryptMessage(pCtx context.Context, pNetMsg net_message.IMessage) (asymmetric.IPubKey, []byte, error) {
	resp, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleMessageDecryptTemplate, p.fHost),
		pNetMsg.ToString(),
	)
	if err != nil {
		return nil, nil, utils.MergeErrors(ErrRequest, err)
	}

	var result hle_settings.SContainer
	if err := encoding.DeserializeJSON(resp, &result); err != nil {
		return nil, nil, utils.MergeErrors(ErrDecodeResponse, err)
	}

	pubKey := asymmetric.LoadRSAPubKey(result.FPublicKey)
	if pubKey == nil {
		return nil, nil, utils.MergeErrors(ErrInvalidResponse, errors.New("decode public key"))
	}

	data := encoding.HexDecode(result.FHexData)
	if data == nil {
		return nil, nil, utils.MergeErrors(ErrInvalidResponse, errors.New("decode data"))
	}

	return pubKey, data, nil
}

func (p *sRequester) GetPubKey(pCtx context.Context) (asymmetric.IPubKey, error) {
	res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleServicePubKeyTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return nil, utils.MergeErrors(ErrRequest, err)
	}

	pubKey := asymmetric.LoadRSAPubKey(string(res))
	if pubKey == nil {
		return nil, utils.MergeErrors(ErrDecodeResponse, errors.New("got invalid public key"))
	}

	return pubKey, nil
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
		return nil, utils.MergeErrors(ErrRequest, err)
	}

	cfgSettings := new(config.SConfigSettings)
	if err := encoding.DeserializeJSON([]byte(res), cfgSettings); err != nil {
		return nil, utils.MergeErrors(ErrDecodeResponse, err)
	}

	return cfgSettings, nil
}
