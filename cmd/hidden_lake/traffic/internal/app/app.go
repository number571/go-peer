package app

import (
	"errors"
	"net/http"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/settings"
	"github.com/number571/go-peer/pkg/closer"
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
	fServiceHTTP *http.Server
}

func NewApp(
	cfg config.IConfig,
) types.IApp {
	db := database.NewKeyValueDB(
		hlt_settings.CPathDB,
		database.NewSettings(&database.SSettings{
			FLimitMessages: hlt_settings.CLimitMessages,
			FMessageSize:   hlt_settings.CMessageSize,
			FWorkSize:      hlt_settings.CWorkSize,
		}),
	)
	return &sApp{
		fConfig:      cfg,
		fDatabase:    db,
		fServiceHTTP: initServiceHTTP(cfg, db),
	}
}

func (app *sApp) Run() error {
	res := make(chan error)

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
		app.fServiceHTTP,
		app.fDatabase,
	})
}
