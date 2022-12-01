package main

import (
	"fmt"

	"github.com/number571/go-peer/cmd/hlm/app"
	"github.com/number571/go-peer/cmd/hlm/config"
	"github.com/number571/go-peer/cmd/hlm/database"
	"github.com/number571/go-peer/cmd/hls/hlc"
	"github.com/number571/go-peer/modules/filesystem"

	hlm_settings "github.com/number571/go-peer/cmd/hlm/settings"
)

func initValues() error {
	cfg, err := getConfig()
	if err != nil {
		return err
	}

	levelDB := database.NewKeyValueDB(hlm_settings.CPathDB, "")
	if levelDB == nil {
		return fmt.Errorf("error: create/open database")
	}

	client := hlc.NewClient(
		hlc.NewRequester(fmt.Sprintf("http://%s", cfg.Connection())),
	)

	if _, err := client.PubKey(); err != nil {
		return err
	}

	gApp = app.NewApp(cfg, client, levelDB)
	return nil
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
