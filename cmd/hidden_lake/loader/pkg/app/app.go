package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/number571/go-peer/cmd/hidden_lake/loader/internal/config"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/loader/pkg/settings"
	"github.com/number571/go-peer/internal/interrupt"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	std_logger "github.com/number571/go-peer/internal/logger/std"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/state"
	"github.com/number571/go-peer/pkg/types"
)

var (
	_ types.IRunner = &sApp{}
)

type sApp struct {
	fState state.IState
	fMutex sync.Mutex

	fConfig config.IConfig

	fHTTPLogger logger.ILogger
	fStdfLogger logger.ILogger

	fServiceHTTP  *http.Server
	fServicePPROF *http.Server
}

func NewApp(
	pCfg config.IConfig,
	pPathTo string,
) types.IRunner {
	logging := pCfg.GetLogging()

	var (
		httpLogger = std_logger.NewStdLogger(logging, http_logger.GetLogFunc())
		stdfLogger = std_logger.NewStdLogger(logging, std_logger.GetLogFunc())
	)

	return &sApp{
		fState:      state.NewBoolState(),
		fConfig:     pCfg,
		fHTTPLogger: httpLogger,
		fStdfLogger: stdfLogger,
	}
}

func (p *sApp) Run(pCtx context.Context) error {
	enableFunc := func() error {
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

	p.initServiceHTTP()
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
		if p.fConfig.GetAddress().GetHTTP() == "" {
			return
		}

		err := p.fServiceHTTP.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			chErr <- err
			return
		}
	}()

	select {
	case <-pCtx.Done():
		return nil
	case err := <-chErr:
		return fmt.Errorf("got run error: %w", err)
	}
}

func (p *sApp) stop() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	err := interrupt.CloseAll([]types.ICloser{
		p.fServiceHTTP,
	})
	if err != nil {
		return fmt.Errorf("close/stop all: %w", err)
	}

	return nil
}
