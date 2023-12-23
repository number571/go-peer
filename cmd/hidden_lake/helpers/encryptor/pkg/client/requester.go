package client

import (
	"errors"
	"fmt"
	"net/http"

	hle_settings "github.com/number571/go-peer/cmd/hidden_lake/helpers/encryptor/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	net_message "github.com/number571/go-peer/pkg/network/message"
)

const (
	cHandleIndexTemplate   = "%s" + hle_settings.CHandleIndexPath
	cHandleEncryptTemplate = "%s" + hle_settings.CHandleEncryptPath
	cHandleDecryptTemplate = "%s" + hle_settings.CHandleDecryptPath
	cHandlePubKeyTemplate  = "%s" + hle_settings.CHandlePubKeyPath
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

func (p *sRequester) GetIndex() (string, error) {
	res, err := api.Request(
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleIndexTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("get index (requester): %w", err)
	}

	result := string(res)
	if result != hle_settings.CTitlePattern {
		return "", errors.New("incorrect title pattern")
	}

	return result, nil
}

func (p *sRequester) EncryptMessage(pPubKey asymmetric.IPubKey, pData []byte) (net_message.IMessage, error) {
	resp, err := api.Request(
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleEncryptTemplate, p.fHost),
		hle_settings.SContainer{
			FPublicKey: pPubKey.ToString(),
			FHexData:   encoding.HexEncode(pData),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("encrypt message (requester): %w", err)
	}

	msg, err := net_message.LoadMessage(p.fParams, string(resp))
	if err != nil {
		return nil, fmt.Errorf("load message (requester): %w", err)
	}

	return msg, nil
}

func (p *sRequester) DecryptMessage(pNetMsg net_message.IMessage) (asymmetric.IPubKey, []byte, error) {
	resp, err := api.Request(
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleDecryptTemplate, p.fHost),
		pNetMsg.ToString(),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("decrypt message (requester): %w", err)
	}

	var result hle_settings.SContainer
	if err := encoding.DeserializeJSON(resp, &result); err != nil {
		return nil, nil, fmt.Errorf("decode response (requester): %w", err)
	}

	pubKey := asymmetric.LoadRSAPubKey(result.FPublicKey)
	if pubKey == nil {
		return nil, nil, errors.New("decode public key (requester)")
	}

	data := encoding.HexDecode(result.FHexData)
	if data == nil {
		return nil, nil, errors.New("decode data (requester)")
	}

	return pubKey, data, nil
}

func (p *sRequester) GetPubKey() (asymmetric.IPubKey, error) {
	res, err := api.Request(
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandlePubKeyTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("get public key (requester): %w", err)
	}

	pubKey := asymmetric.LoadRSAPubKey(string(res))
	if pubKey == nil {
		return nil, errors.New("got invalid public key")
	}

	return pubKey, nil
}
