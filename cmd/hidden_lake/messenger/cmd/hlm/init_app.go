package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/app"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/database"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/types"

	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/client"
)

func initApp() (types.ICommand, error) {
	cfg, err := initConfig()
	if err != nil {
		return nil, err
	}

	stg, err := initCryptoStorage(cfg)
	if err != nil {
		return nil, err
	}

	hlsClient := hls_client.NewClient(
		hls_client.NewBuilder(),
		hls_client.NewRequester(
			fmt.Sprintf("http://%s", cfg.GetConnection().GetService()),
			&http.Client{Timeout: time.Minute},
		),
	)

	hltClient := hlt_client.NewClient(
		hlt_client.NewBuilder(),
		hlt_client.NewRequester(
			fmt.Sprintf("http://%s", cfg.GetConnection().GetTraffic()),
			&http.Client{Timeout: time.Minute},
			message.NewParams(hls_settings.CMessageSize, hls_settings.CWorkSize),
		),
	)

	wDB := database.NewWrapperDB()
	return app.NewApp(cfg, stg, wDB, hlsClient, hltClient), nil
}
