package main

import (
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/app"
	internal_logger "github.com/number571/go-peer/internal/logger"
	"github.com/number571/go-peer/pkg/types"
)

func initApp() (types.IApp, error) {
	cfg, err := initConfig()
	if err != nil {
		return nil, err
	}

	db := initDatabase()
	logger := internal_logger.DefaultLogger(cfg.Logging())
	connKeeper := initConnKeeper(cfg, db, logger)

	return app.NewApp(cfg, db, connKeeper), nil
}
