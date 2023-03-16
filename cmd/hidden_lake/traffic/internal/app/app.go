package app

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	"github.com/number571/go-peer/pkg/network/conn_keeper"
	"github.com/number571/go-peer/pkg/types"

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

	fConfig config.IConfig

	fDatabase    database.IKeyValueDB
	fConnKeeper  conn_keeper.IConnKeeper
	fServiceHTTP *http.Server
}

func NewApp(
	cfg config.IConfig,
) types.ICommand {
	return &sApp{
		fConfig: cfg,
	}
}

func (app *sApp) Run() error {
	app.fMutex.Lock()
	defer app.fMutex.Unlock()

	if app.fIsRun {
		return errors.New("application already running")
	}
	app.fIsRun = true

	logger := internal_logger.StdLogger(app.fConfig.GetLogging())

	app.fDatabase = initDatabase()
	app.fConnKeeper = initConnKeeper(app.fConfig, app.fDatabase, logger)
	app.fServiceHTTP = initServiceHTTP(app.fConfig, app.fConnKeeper, app.fDatabase)

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
	app.fMutex.Lock()
	defer app.fMutex.Unlock()

	if !app.fIsRun {
		return errors.New("application already stopped or not started")
	}
	app.fIsRun = false

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
