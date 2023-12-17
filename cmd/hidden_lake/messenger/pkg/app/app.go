package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/database"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/go-peer/pkg/utils"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/settings"
	"github.com/number571/go-peer/internal/interrupt"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	std_logger "github.com/number571/go-peer/internal/logger/std"
)

var (
	_ types.IRunner = &sApp{}
)

type sApp struct {
	fIsRun bool
	fMutex sync.Mutex

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
		fConfig:     pCfg,
		fPathTo:     pPathTo,
		fHTTPLogger: httpLogger,
		fStdfLogger: stdfLogger,
	}
}

func (p *sApp) Run(pCtx context.Context) error {
	err := func() error {
		p.fMutex.Lock()
		defer p.fMutex.Unlock()

		if p.fIsRun {
			return errors.New("application already running")
		}

		if err := p.initDatabase(); err != nil {
			return fmt.Errorf("init database: %w", err)
		}

		p.fStdfLogger.PushInfo(fmt.Sprintf("%s is running...", pkg_settings.CServiceName))
		p.fIsRun = true

		return nil
	}()
	if err != nil {
		return err
	}

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
		return p.stop()
	case err := <-chErr:
		return utils.MergeErrors(
			fmt.Errorf("got run error: %w", err),
			p.stop(),
		)
	}
}

func (p *sApp) stop() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	p.fStdfLogger.PushInfo(fmt.Sprintf("%s is shutting down...", pkg_settings.CServiceName))
	p.fIsRun = false

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
