package main

import (
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/app"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"

	"github.com/number571/go-peer/pkg/types"
)

func initApp() (types.ICommand, error) {
	cfg, err := config.InitConfig(settings.CPathCFG, nil)
	if err != nil {
		return nil, err
	}

	return app.NewApp(cfg, "."), nil
}
