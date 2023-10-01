package app

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/app/state"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/types"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/settings"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	std_logger "github.com/number571/go-peer/internal/logger/std"
	pkg_errors "github.com/number571/go-peer/pkg/errors"
)

const (
	cInitStart = time.Second * 3
)

var (
	_ types.ICommand = &sApp{}
)

type sApp struct {
	fIsRun bool
	fMutex sync.Mutex

	fWrapper      config.IWrapper
	fStateManager state.IStateManager

	fIntServiceHTTP *http.Server
	fIncServiceHTTP *http.Server
	fServicePPROF   *http.Server

	fHTTPLogger logger.ILogger
	fStdfLogger logger.ILogger
}

func NewApp(
	pCfg config.IConfig,
	pPathTo string,
) types.ICommand {
	httpLogger := std_logger.NewStdLogger(pCfg.GetLogging(), http_logger.GetLogFunc())
	stdfLogger := std_logger.NewStdLogger(pCfg.GetLogging(), std_logger.GetLogFunc())

	return &sApp{
		fWrapper:      config.NewWrapper(pCfg),
		fStateManager: state.NewStateManager(pCfg, pPathTo),
		fHTTPLogger:   httpLogger,
		fStdfLogger:   stdfLogger,
	}
}

func (p *sApp) Run() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if p.fIsRun {
		return pkg_errors.NewError("application already running")
	}
	p.fIsRun = true

	p.initIncomingServiceHTTP()
	p.initInterfaceServiceHTTP()
	p.initServicePPROF()

	res := make(chan error)

	go func() {
		if p.fWrapper.GetConfig().GetAddress().GetPPROF() == "" {
			return
		}

		err := p.fServicePPROF.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			res <- err
			return
		}
	}()

	go func() {
		err := p.fIntServiceHTTP.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			res <- err
			return
		}
	}()

	go func() {
		err := p.fIncServiceHTTP.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			res <- err
			return
		}
	}()

	select {
	case err := <-res:
		return pkg_errors.AppendError(pkg_errors.WrapError(err, "got run error"), p.Stop())
	case <-time.After(cInitStart):
		p.fStdfLogger.PushInfo(fmt.Sprintf("%s is running...", pkg_settings.CServiceName))
		return nil
	}
}

func (p *sApp) Stop() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if !p.fIsRun {
		return pkg_errors.NewError("application already stopped or not started")
	}
	p.fIsRun = false
	p.fStdfLogger.PushInfo(fmt.Sprintf("%s is shutting down...", pkg_settings.CServiceName))

	// state may be already closed by HLS
	_ = p.fStateManager.CloseState()

	err := types.CloseAll([]types.ICloser{
		p.fIntServiceHTTP,
		p.fIncServiceHTTP,
		p.fServicePPROF,
		p.fStateManager.GetWrapperDB(),
	})
	if err != nil {
		return pkg_errors.WrapError(err, "close/stop all")
	}
	return nil
}
