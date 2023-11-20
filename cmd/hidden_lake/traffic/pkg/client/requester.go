package client

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/encoding"
	net_message "github.com/number571/go-peer/pkg/network/message"
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
	resp, err := api.Request(
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(pkg_settings.CHandleIndexTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("get index (requester): %w", err)
	}

	result := string(resp)
	if result != pkg_settings.CTitlePattern {
		return "", errors.New("incorrect title pattern")
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
		return nil, fmt.Errorf("get hashes (requester): %w", err)
	}

	var hashes []string
	if err := encoding.Deserialize([]byte(resp), &hashes); err != nil {
		return nil, fmt.Errorf("deserialize hashes (requeser): %w", err)
	}

	return hashes, nil
}

func (p *sRequester) GetMessage(pHash string) (net_message.IMessage, error) {
	resp, err := api.Request(
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(pkg_settings.CHandleMessageTemplate+"?hash=%s", p.fHost, pHash),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("get message (requester): %w", err)
	}

	msg := net_message.LoadMessage(p.fParams, string(resp))
	if msg == nil {
		return nil, errors.New("load message")
	}

	if !bytes.Equal(msg.GetHash(), encoding.HexDecode(pHash)) {
		return nil, errors.New("got invalid hash")
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
		return fmt.Errorf("put message (requester): %w", err)
	}
	return nil
}
