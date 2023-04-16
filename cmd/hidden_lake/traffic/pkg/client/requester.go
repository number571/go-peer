package client

import (
	"fmt"
	"net/http"
	"strings"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/encoding"
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
		return "", err
	}

	if resp != pkg_settings.CTitlePattern {
		return "", fmt.Errorf("incorrect title pattern")
	}
	return resp, nil
}

func (p *sRequester) GetHashes() ([]string, error) {
	resp, err := api.Request(
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(pkg_settings.CHandleHashesTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return nil, err
	}
	return strings.Split(resp, ";"), nil
}

func (p *sRequester) GetMessage(pRequest *pkg_settings.SLoadRequest) (message.IMessage, error) {
	resp, err := api.Request(
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(pkg_settings.CHandleMessageTemplate+"?hash=%s", p.fHost, pRequest.FHash),
		nil,
	)
	if err != nil {
		return nil, err
	}

	msg := message.LoadMessage(
		message.NewSettings(&message.SSettings{
			FWorkSize:    p.fParams.GetWorkSize(),
			FMessageSize: p.fParams.GetMessageSize(),
		}),
		encoding.HexDecode(resp),
	)
	if msg == nil {
		return nil, fmt.Errorf("message is nil")
	}

	return msg, nil
}

func (p *sRequester) PutMessage(pRequest *pkg_settings.SPushRequest) error {
	_, err := api.Request(
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(pkg_settings.CHandleMessageTemplate, p.fHost),
		pRequest,
	)
	return err
}
