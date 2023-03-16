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
	fParams message.IParams
}

func NewRequester(host string, client *http.Client, params message.IParams) IRequester {
	return &sRequester{
		fHost:   host,
		fClient: client,
		fParams: params,
	}
}

func (r *sRequester) GetIndex() (string, error) {
	resp, err := api.Request(
		r.fClient,
		http.MethodGet,
		fmt.Sprintf(pkg_settings.CHandleIndexTemplate, r.fHost),
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

func (r *sRequester) GetHashes() ([]string, error) {
	resp, err := api.Request(
		r.fClient,
		http.MethodGet,
		fmt.Sprintf(pkg_settings.CHandleHashesTemplate, r.fHost),
		nil,
	)
	if err != nil {
		return nil, err
	}
	return strings.Split(resp, ";"), nil
}

func (r *sRequester) GetMessage(request *pkg_settings.SLoadRequest) (message.IMessage, error) {
	resp, err := api.Request(
		r.fClient,
		http.MethodGet,
		fmt.Sprintf(pkg_settings.CHandleMessageTemplate+"?hash=%s", r.fHost, request.FHash),
		nil,
	)
	if err != nil {
		return nil, err
	}

	msg := message.LoadMessage(
		encoding.HexDecode(resp),
		message.NewParams(
			r.fParams.GetMessageSize(),
			r.fParams.GetWorkSize(),
		),
	)
	if msg == nil {
		return nil, fmt.Errorf("message is nil")
	}

	return msg, nil
}

func (r *sRequester) PutMessage(request *pkg_settings.SPushRequest) error {
	_, err := api.Request(
		r.fClient,
		http.MethodPost,
		fmt.Sprintf(pkg_settings.CHandleMessageTemplate, r.fHost),
		request,
	)
	return err
}
