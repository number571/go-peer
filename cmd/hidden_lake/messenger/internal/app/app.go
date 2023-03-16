package app

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/app/state"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/database"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/types"

	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/client"
)

const (
	initStart = time.Second * 3
)

var (
	_ types.ICommand = &sApp{}
)

type sApp struct {
	fState          state.IState
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
		fIntServiceHTTP: initInterfaceServiceHTTP(cfg, state),
		fIncServiceHTTP: initIncomingServiceHTTP(cfg, state),
	}
}

func (app *sApp) Run() error {
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
		return nil
	}
}

func (app *sApp) Stop() error {
	_ = app.fState.ClearActiveState()

	lastErr := types.CloseAll([]types.ICloser{
		app.fIntServiceHTTP,
		app.fIncServiceHTTP,
	})

	db := app.fState.GetWrapperDB().Get()
	if db != nil {
		if err := db.Close(); err != nil {
			lastErr = err
		}
	}

	return lastErr
}
