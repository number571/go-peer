package app

import (
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
	_ types.IApp = &sApp{}
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
) types.IApp {
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

func (p *sApp) Run() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if p.fIsRun {
		return errors.New("application already running")
	}
	p.fIsRun = true

	p.initServiceHTTP()
	p.initServicePPROF()

	res := make(chan error)

	go func() {
		if p.fConfig.GetAddress().GetPPROF() == "" {
			return
		}

		err := p.fServicePPROF.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			res <- err
			return
		}
	}()

	go func() {
		if p.fConfig.GetAddress().GetHTTP() == "" {
			return
		}

		err := p.fServiceHTTP.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			res <- err
			return
		}
	}()

	select {
	case err := <-res:
		resErr := fmt.Errorf("got run error: %w", err)
		return utils.MergeErrors(resErr, p.Stop())
	case <-time.After(cInitStart):
		p.fStdfLogger.PushInfo(fmt.Sprintf("%s is running...", settings.CServiceName))
		return nil
	}
}

func (p *sApp) Stop() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if !p.fIsRun {
		return errors.New("application already stopped or not started")
	}
	p.fIsRun = false
	p.fStdfLogger.PushInfo(fmt.Sprintf("%s is shutting down...", settings.CServiceName))

	err := interrupt.CloseAll([]types.ICloser{
		p.fServiceHTTP,
	})
	if err != nil {
		return fmt.Errorf("close/stop all: %w", err)
	}

	return nil
}
