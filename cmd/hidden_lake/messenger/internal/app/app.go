package app

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/app/state"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/database"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/settings"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/types"

	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/client"
	internal_logger "github.com/number571/go-peer/internal/logger"
)

const (
	initStart = time.Second * 3
)

var (
	_ types.ICommand = &sApp{}
)

type sApp struct {
	fIsRun bool
	fMutex sync.Mutex

	fState          state.IState
	fLogger         logger.ILogger
	fIntServiceHTTP *http.Server
	fIncServiceHTTP *http.Server
}

func NewApp(
	cfg config.IConfig,
) types.ICommand {
	stg, err := initCryptoStorage(cfg)
	if err != nil {
		panic(err)
	}

	state := state.NewState(
		stg,
		database.NewWrapperDB(),
		hls_client.NewClient(
			hls_client.NewBuilder(),
			hls_client.NewRequester(
				fmt.Sprintf("http://%s", cfg.GetConnection().GetService()),
				&http.Client{Timeout: time.Minute},
			),
		),
		hlt_client.NewClient(
			hlt_client.NewBuilder(),
			hlt_client.NewRequester(
				fmt.Sprintf("http://%s", cfg.GetConnection().GetTraffic()),
				&http.Client{Timeout: time.Minute},
				message.NewParams(hls_settings.CMessageSize, hls_settings.CWorkSize),
			),
		),
	)

	return &sApp{
		fState:          state,
		fLogger:         internal_logger.StdLogger(cfg.GetLogging()),
		fIntServiceHTTP: initInterfaceServiceHTTP(cfg, state),
		fIncServiceHTTP: initIncomingServiceHTTP(cfg, state),
	}
}

func (app *sApp) Run() error {
	app.fMutex.Lock()
	defer app.fMutex.Unlock()

	if app.fIsRun {
		return errors.New("application already running")
	}
	app.fIsRun = true

	res := make(chan error)

	go func() {
		err := app.fIntServiceHTTP.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			res <- err
			return
		}
	}()

	go func() {
		err := app.fIncServiceHTTP.ListenAndServe()
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
		app.fLogger.PushInfo(fmt.Sprintf("%s is running...", settings.CServiceName))
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
	app.fLogger.PushInfo(fmt.Sprintf("%s is shutting down...", settings.CServiceName))

	// state may be already closed
	_ = app.fState.ClearActiveState()

	return types.CloseAll([]types.ICloser{
		app.fIntServiceHTTP,
		app.fIncServiceHTTP,
		app.fState.GetWrapperDB(),
	})
}
