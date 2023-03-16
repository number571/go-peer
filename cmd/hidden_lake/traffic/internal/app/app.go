package app

import (
	"errors"
	"net/http"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	"github.com/number571/go-peer/pkg/network/conn_keeper"
	"github.com/number571/go-peer/pkg/types"
)

const (
	initStart = 3 * time.Second
)

var (
	_ types.ICommand = &sApp{}
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
) types.ICommand {
	return &sApp{
		fConfig:      cfg,
		fDatabase:    db,
		fConnKeeper:  connKeeper,
		fServiceHTTP: initServiceHTTP(cfg, connKeeper, db),
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
		if app.fConfig.GetAddress() == "" {
			return
		}

		err := app.fServiceHTTP.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			res <- err
			return
		}
	}()

	select {
	case err := <-res:
		app.Stop()
		return err
	case <-time.After(initStart):
		return nil
	}
}

func (app *sApp) Stop() error {
	lastErr := types.StopAll([]types.ICommand{
		app.fConnKeeper,
		app.fConnKeeper.GetNetworkNode(),
	})

	err := types.CloseAll([]types.ICloser{
		app.fServiceHTTP,
		app.fDatabase,
	})
	if err != nil {
		lastErr = err
	}

	return lastErr
}
