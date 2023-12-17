package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/loader/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/loader/pkg/settings"
	"github.com/number571/go-peer/internal/interrupt"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	std_logger "github.com/number571/go-peer/internal/logger/std"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/go-peer/pkg/utils"
)

const (
	cInitStart = 3 * time.Second
)

var (
	_ types.IRunner = &sApp{}
)

type sApp struct {
	fIsRun bool
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
		fConfig:     pCfg,
		fHTTPLogger: httpLogger,
		fStdfLogger: stdfLogger,
	}
}

func (p *sApp) Run(pCtx context.Context) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if p.fIsRun {
		return errors.New("application already running")
	}

	p.initServiceHTTP()
	p.initServicePPROF()

	p.fIsRun = true
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

	deadline, ok := pCtx.Deadline()
	if !ok {
		// set default value
		deadline = time.Now().Add(cInitStart)
	}

	dlCtx, cancel := context.WithDeadline(pCtx, deadline)
	defer cancel()

	select {
	case err := <-chErr:
		return utils.MergeErrors(
			fmt.Errorf("got run error: %w", err),
			p.stop(),
		)

	case <-dlCtx.Done():
		p.fStdfLogger.PushInfo(fmt.Sprintf("%s is running...", settings.CServiceName))
		p.fIsRun = true

		go func() {
			<-pCtx.Done()

			p.fMutex.Lock()
			_ = p.stop()
			p.fMutex.Unlock()
		}()

		return nil
	}
}

func (p *sApp) stop() error {
	p.fStdfLogger.PushInfo(fmt.Sprintf("%s is shutting down...", settings.CServiceName))
	p.fIsRun = false

	err := interrupt.CloseAll([]types.ICloser{
		p.fServiceHTTP,
	})
	if err != nil {
		return fmt.Errorf("close/stop all: %w", err)
	}

	return nil
}
