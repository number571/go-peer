package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/database"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/state"
	"github.com/number571/go-peer/pkg/types"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/settings"
	"github.com/number571/go-peer/internal/interrupt"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	std_logger "github.com/number571/go-peer/internal/logger/std"
)

var (
	_ types.IRunner = &sApp{}
)

type sApp struct {
	fState state.IState

	fConfig config.IConfig
	fPathTo string

	fDatabase       database.IKVDatabase
	fIntServiceHTTP *http.Server
	fIncServiceHTTP *http.Server
	fServicePPROF   *http.Server

	fHTTPLogger logger.ILogger
	fStdfLogger logger.ILogger
}

func NewApp(
	pCfg config.IConfig,
	pPathTo string,
) types.IRunner {
	httpLogger := std_logger.NewStdLogger(pCfg.GetLogging(), http_logger.GetLogFunc())
	stdfLogger := std_logger.NewStdLogger(pCfg.GetLogging(), std_logger.GetLogFunc())

	return &sApp{
		fState:      state.NewBoolState(),
		fConfig:     pCfg,
		fPathTo:     pPathTo,
		fHTTPLogger: httpLogger,
		fStdfLogger: stdfLogger,
	}
}

func (p *sApp) Run(pCtx context.Context) error {
	enableFunc := func() error {
		if err := p.initDatabase(); err != nil {
			return fmt.Errorf("init database: %w", err)
		}
		p.fStdfLogger.PushInfo(fmt.Sprintf("%s is running...", pkg_settings.CServiceName))
		return nil
	}
	if err := p.fState.Enable(enableFunc); err != nil {
		return fmt.Errorf("application running error: %w", err)
	}

	defer func() {
		disableFunc := func() error {
			p.fStdfLogger.PushInfo(fmt.Sprintf("%s is shutting down...", pkg_settings.CServiceName))
			return p.stop()
		}
		_ = p.fState.Disable(disableFunc)
	}()

	p.initIncomingServiceHTTP()
	p.initInterfaceServiceHTTP()
	p.initServicePPROF()

	chErr := make(chan error)

	go func() {
		if p.fConfig.GetAddress().GetPPROF() == "" {
			return
		}

		err := p.fServicePPROF.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			chErr <- err
			return
		}
	}()

	go func() {
		err := p.fIntServiceHTTP.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			chErr <- err
			return
		}
	}()

	go func() {
		err := p.fIncServiceHTTP.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			chErr <- err
			return
		}
	}()

	select {
	case <-pCtx.Done():
		return pCtx.Err()
	case err := <-chErr:
		return fmt.Errorf("got run error: %w", err)
	}
}

func (p *sApp) stop() error {
	err := interrupt.CloseAll([]types.ICloser{
		p.fIntServiceHTTP,
		p.fIncServiceHTTP,
		p.fServicePPROF,
		p.fDatabase,
	})
	if err != nil {
		return fmt.Errorf("close/stop all: %w", err)
	}
	return nil
}
