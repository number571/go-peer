package client

import (
	"errors"
	"fmt"
	"net/http"

	hll_settings "github.com/number571/go-peer/cmd/hidden_lake/helpers/loader/pkg/settings"
	"github.com/number571/go-peer/internal/api"
)

const (
	cHandleIndexTemplate    = "%s" + hll_settings.CHandleIndexPath
	cHandleTransferTemplate = "%s" + hll_settings.CHandleTransferPath
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
	if result != hll_settings.CTitlePattern {
		return "", errors.New("incorrect title pattern")
	}

	return result, nil
}

func (p *sRequester) RunTransfer() error {
	_, err := api.Request(
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleTransferTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return fmt.Errorf("run loader (requester): %w", err)
	}
	return nil
}

func (p *sRequester) StopTransfer() error {
	_, err := api.Request(
		p.fClient,
		http.MethodDelete,
		fmt.Sprintf(cHandleTransferTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return fmt.Errorf("stop loader (requester): %w", err)
	}
	return nil
}
