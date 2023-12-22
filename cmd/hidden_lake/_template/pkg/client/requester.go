package client

import (
	"errors"
	"fmt"
	"net/http"

	hl_t_settings "github.com/number571/go-peer/cmd/hidden_lake/_template/pkg/settings"
	"github.com/number571/go-peer/internal/api"
)

const (
	cHandleIndexTemplate = "%s" + hl_t_settings.CHandleIndexPath
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
	if result != hl_t_settings.CTitlePattern {
		return "", errors.New("incorrect title pattern")
	}

	return result, nil
}
