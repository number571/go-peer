package client

import (
	"bytes"
	"fmt"
	"net/http"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/errors"
)

var (
	_ IRequester = &sRequester{}
)

type sRequester struct {
	fHost   string
	fClient *http.Client
	fParams message.ISettings
}

func NewRequester(pHost string, pClient *http.Client, pParams message.ISettings) IRequester {
	return &sRequester{
		fHost:   pHost,
		fClient: pClient,
		fParams: pParams,
	}
}

func (p *sRequester) GetIndex() (string, error) {
	resp, err := api.Request(
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(pkg_settings.CHandleIndexTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return "", errors.WrapError(err, "get index (requester)")
	}

	result := string(resp)
	if result != pkg_settings.CTitlePattern {
		return "", errors.NewError("incorrect title pattern")
	}

	return result, nil
}

func (p *sRequester) GetHashes() ([]string, error) {
	resp, err := api.Request(
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(pkg_settings.CHandleHashesTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return nil, errors.WrapError(err, "get hashes (requester)")
	}

	var hashes []string
	if err := encoding.Deserialize([]byte(resp), &hashes); err != nil {
		return nil, errors.WrapError(err, "deserialize hashes (requeser)")
	}

	return hashes, nil
}

func (p *sRequester) GetMessage(pHash string) (message.IMessage, error) {
	resp, err := api.Request(
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(pkg_settings.CHandleMessageTemplate+"?hash=%s", p.fHost, pHash),
		nil,
	)
	if err != nil {
		return nil, errors.WrapError(err, "get message (requester)")
	}

	msg := message.LoadMessage(p.fParams, string(resp))
	if msg == nil {
		return nil, errors.NewError("load message")
	}

	if !bytes.Equal(msg.GetBody().GetHash(), encoding.HexDecode(pHash)) {
		return nil, errors.NewError("got invalid hash")
	}
	return msg, nil
}

func (p *sRequester) PutMessage(pRequest string) error {
	_, err := api.Request(
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(pkg_settings.CHandleMessageTemplate, p.fHost),
		pRequest,
	)
	if err != nil {
		return errors.WrapError(err, "put message (requester)")
	}
	return nil
}
