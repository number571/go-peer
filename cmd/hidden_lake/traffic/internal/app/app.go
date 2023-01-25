package app

import (
	"errors"
	"net/http"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	"github.com/number571/go-peer/pkg/closer"
	"github.com/number571/go-peer/pkg/network/conn_keeper"
	"github.com/number571/go-peer/pkg/types"
)

const (
	initStart = 3 * time.Second
)

var (
	_ types.IApp = &sApp{}
)

type sApp struct {
	fConfig      config.IConfig
	fDatabase    database.IKeyValueDB
	fConnKeeper  conn_keeper.IConnKeeper
	fServiceHTTP *http.Server
}

func NewApp(
	cfg config.IConfig,
	db database.IKeyValueDB,
	connKeeper conn_keeper.IConnKeeper,
) types.IApp {
	return &sApp{
		fConfig:      cfg,
		fDatabase:    db,
		fConnKeeper:  connKeeper,
		fServiceHTTP: initServiceHTTP(cfg, db),
	}
}

func (app *sApp) Run() error {
	res := make(chan error)

	go func() {
		if err := app.fConnKeeper.Run(); err != nil {
			res <- err
			return
		}
	}()

	go func() {
		err := app.fServiceHTTP.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			res <- err
			return
		}
	}()

	select {
	case err := <-res:
		app.Close()
		return err
	case <-time.After(initStart):
		return nil
	}
}

func (app *sApp) Close() error {
	return closer.CloseAll([]types.ICloser{
		app.fConnKeeper.Network(),
		app.fConnKeeper,
		app.fServiceHTTP,
		app.fDatabase,
	})
}
