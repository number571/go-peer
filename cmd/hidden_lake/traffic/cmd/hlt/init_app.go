package main

import (
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/app"
	"github.com/number571/go-peer/pkg/types"
)

func initApp() (types.IApp, error) {
	cfg, err := initConfig()
	if err != nil {
		return nil, err
	}

	db := initDatabase()
	connKeeper := initConnKeeper(cfg, db)

	return app.NewApp(cfg, db, connKeeper), nil
}
