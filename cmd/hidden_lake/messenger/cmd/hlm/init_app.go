package main

import (
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/app"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/settings"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/types"
)

func initApp() (types.ICommand, error) {
	cfg, err := config.InitConfig(settings.CPathCFG, nil)
	if err != nil {
		return nil, errors.WrapError(err, "init config")
	}

	return app.NewApp(cfg, "."), nil
}
