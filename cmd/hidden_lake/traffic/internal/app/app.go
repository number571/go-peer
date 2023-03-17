package app

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network/conn_keeper"
	"github.com/number571/go-peer/pkg/types"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	internal_logger "github.com/number571/go-peer/internal/logger"
)

const (
	initStart = 3 * time.Second
)

var (
	_ types.ICommand = &sApp{}
)

type sApp struct {
	fIsRun bool
	fMutex sync.Mutex

	fConfig      config.IConfig
	fWrapperDB   database.IWrapperDB
	fLogger      logger.ILogger
	fConnKeeper  conn_keeper.IConnKeeper
	fServiceHTTP *http.Server
}

func NewApp(
	cfg config.IConfig,
) types.ICommand {
	wDB := database.NewWrapperDB()
	logger := internal_logger.StdLogger(cfg.GetLogging())
	connKeeper := initConnKeeper(cfg, wDB, logger)
	return &sApp{
		fConfig:      cfg,
		fWrapperDB:   wDB,
		fLogger:      logger,
		fConnKeeper:  connKeeper,
		fServiceHTTP: initServiceHTTP(cfg, connKeeper, wDB),
	}
}

func (app *sApp) Run() error {
	app.fMutex.Lock()
	defer app.fMutex.Unlock()

	if app.fIsRun {
		return errors.New("application already running")
	}
	app.fIsRun = true

	app.fWrapperDB.Set(initDatabase())
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
		app.fLogger.PushInfo(fmt.Sprintf("%s is running...", pkg_settings.CServiceName))
		return nil
	}
}

func (app *sApp) Stop() error {
	app.fMutex.Lock()
	defer app.fMutex.Unlock()

	if !app.fIsRun {
		return errors.New("application already stopped or not started")
	}
	app.fIsRun = false
	app.fLogger.PushInfo(fmt.Sprintf("%s is shutting down...", pkg_settings.CServiceName))

	lastErr := types.StopAll([]types.ICommand{
		app.fConnKeeper,
		app.fConnKeeper.GetNetworkNode(),
	})

	err := types.CloseAll([]types.ICloser{
		app.fServiceHTTP,
		app.fWrapperDB,
	})
	if err != nil {
		lastErr = err
	}

	return lastErr
}
