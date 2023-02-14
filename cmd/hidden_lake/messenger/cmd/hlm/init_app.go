package main

import (
	"fmt"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/app"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/database"
	"github.com/number571/go-peer/pkg/types"

	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/client"
)

func initApp() (types.IApp, error) {
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
		hls_client.NewRequester(fmt.Sprintf("http://%s", cfg.Connection().Service())),
	)

	hltClient := hlt_client.NewClient(
		hlt_client.NewBuilder(),
		hlt_client.NewRequester(fmt.Sprintf("http://%s", cfg.Connection().Traffic())),
	)

	wDB := database.NewWrapperDB()
	return app.NewApp(cfg, stg, wDB, hlsClient, hltClient), nil
}
