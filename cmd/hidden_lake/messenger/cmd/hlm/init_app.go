package main

import (
	"fmt"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/app"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/config"
	"github.com/number571/go-peer/pkg/filesystem"
	"github.com/number571/go-peer/pkg/types"

	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/settings"
	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
)

func initApp() (types.IApp, error) {
	cfg, err := getConfig()
	if err != nil {
		return nil, err
	}

	hlsClient := hls_client.NewClient(
		hls_client.NewRequester(fmt.Sprintf("http://%s", cfg.Connection())),
	)

	if _, err := hlsClient.GetPubKey(); err != nil {
		return nil, err
	}

	return app.NewApp(cfg, hlsClient), nil
}

func getConfig() (config.IConfig, error) {
	if filesystem.OpenFile(hlm_settings.CPathCFG).IsExist() {
		return config.LoadConfig(hlm_settings.CPathCFG)
	}
	initCfg := &config.SConfig{
		FAddress: &config.SAddress{
			FInterface: "localhost:8080",
			FIncoming:  "localhost:8081",
		},
		FConnection: "localhost:9572",
	}
	return config.NewConfig(hlm_settings.CPathCFG, initCfg)
}
