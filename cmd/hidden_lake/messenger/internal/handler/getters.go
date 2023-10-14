package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/utils"
	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
)

type sTemplate struct {
	FLanguage utils.ILanguage
}

func getTemplate(pCfg config.IConfig) *sTemplate {
	return &sTemplate{
		FLanguage: pCfg.GetLanguage(),
	}
}

func getClient(pCfg config.IConfig) hls_client.IClient {
	return hls_client.NewClient(
		hls_client.NewBuilder(),
		hls_client.NewRequester(
			fmt.Sprintf("http://%s", pCfg.GetConnection()),
			&http.Client{Timeout: time.Minute},
		),
	)
}
