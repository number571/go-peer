package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/settings"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/encoding"
)

var (
	_ IRequester = &sRequester{}
)

type sRequester struct {
	fHost string
}

func NewRequester(host string) IRequester {
	return &sRequester{
		fHost: host,
	}
}

func (r *sRequester) GetIndex() (string, error) {
	resp, err := http.Get(
		fmt.Sprintf(pkg_settings.CHandleIndexTemplate, r.fHost),
	)
	if err != nil {
		return "", err
	}

	var response pkg_settings.SResponse
	json.NewDecoder(resp.Body).Decode(&response)

	if response.FReturn != pkg_settings.CErrorNone {
		return "", fmt.Errorf("%s", string(response.FResult))
	}

	if response.FResult != pkg_settings.CTitlePattern {
		return "", fmt.Errorf("incorrect title pattern")
	}

	return response.FResult, nil
}

func (r *sRequester) GetHashes() ([]string, error) {
	resp, err := http.Get(
		fmt.Sprintf(pkg_settings.CHandleHashesTemplate, r.fHost),
	)
	if err != nil {
		return nil, err
	}

	var response pkg_settings.SResponse
	json.NewDecoder(resp.Body).Decode(&response)

	if response.FReturn != pkg_settings.CErrorNone {
		return nil, fmt.Errorf("%s", string(response.FResult))
	}

	return strings.Split(response.FResult, ","), nil
}

func (r *sRequester) GetMessage(request *pkg_settings.SLoadRequest) (message.IMessage, error) {
	resp, err := http.Get(
		fmt.Sprintf(pkg_settings.CHandleMessageTemplate+"?hash=%s", r.fHost, request.FHash),
	)
	if err != nil {
		return nil, err
	}

	var response pkg_settings.SResponse
	json.NewDecoder(resp.Body).Decode(&response)

	if response.FReturn != pkg_settings.CErrorNone {
		return nil, fmt.Errorf("%s", string(response.FResult))
	}

	msg := message.LoadMessage(
		encoding.HexDecode(response.FResult),
		hlt_settings.CMessageSize,
		hlt_settings.CWorkSize,
	)
	if msg == nil {
		return nil, fmt.Errorf("message is nil")
	}

	return msg, nil
}

func (r *sRequester) PutMessage(request *pkg_settings.SPushRequest) error {
	resp, err := http.Post(
		fmt.Sprintf(pkg_settings.CHandleMessageTemplate, r.fHost),
		pkg_settings.CContentType,
		bytes.NewReader(encoding.Serialize(request)),
	)
	if err != nil {
		return err
	}

	var response pkg_settings.SResponse
	json.NewDecoder(resp.Body).Decode(&response)

	if response.FReturn != pkg_settings.CErrorNone {
		return fmt.Errorf("%s", string(response.FResult))
	}

	return nil
}
