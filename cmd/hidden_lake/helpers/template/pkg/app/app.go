package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/template/internal/config"
	hl_t_settings "github.com/number571/go-peer/cmd/hidden_lake/helpers/template/pkg/settings"
	"github.com/number571/go-peer/internal/closer"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	std_logger "github.com/number571/go-peer/internal/logger/std"
	internal_types "github.com/number571/go-peer/internal/types"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/state"
	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/go-peer/pkg/utils"
)

var (
	_ types.IRunner = &sApp{}
)

type sApp struct {
	fState state.IState

	fConfig config.IConfig

	fHTTPLogger logger.ILogger
	fStdfLogger logger.ILogger

	fServiceHTTP  *http.Server
	fServicePPROF *http.Server
}

func NewApp(
	pCfg config.IConfig,
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
	services := []internal_types.IServiceF{
		p.runListenerPPROF,
		p.runListenerHTTP,
	}

	ctx, cancel := context.WithCancel(pCtx)
	defer cancel()

	wg := &sync.WaitGroup{}
	wg.Add(len(services))

	if err := p.fState.Enable(p.enable(ctx)); err != nil {
		return utils.MergeErrors(ErrRunning, err)
	}
	defer func() { _ = p.fState.Disable(p.disable(cancel, wg)) }()

	chErr := make(chan error, len(services))
	for _, f := range services {
		go f(ctx, wg, chErr)
	}

	select {
	case <-pCtx.Done():
		return pCtx.Err()
	case err := <-chErr:
		return utils.MergeErrors(ErrService, err)
	}
}

func (p *sApp) enable(_ context.Context) state.IStateF {
	return func() error {
		p.initServiceHTTP()
		p.initServicePPROF()

		p.fStdfLogger.PushInfo(fmt.Sprintf("%s is running...", hl_t_settings.CServiceName))
		return nil
	}
}

func (p *sApp) disable(pCancel context.CancelFunc, pWg *sync.WaitGroup) state.IStateF {
	return func() error {
		pCancel()
		pWg.Wait() // wait canceled context

		p.fStdfLogger.PushInfo(fmt.Sprintf("%s is shutting down...", hl_t_settings.CServiceName))
		return p.stop()
	}
}

func (p *sApp) runListenerPPROF(pCtx context.Context, wg *sync.WaitGroup, pChErr chan<- error) {
	defer wg.Done()

	if p.fConfig.GetAddress().GetPPROF() == "" {
		return
	}

	go func() {
		err := p.fServicePPROF.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			pChErr <- err
			return
		}
	}()

	<-pCtx.Done()
}

func (p *sApp) runListenerHTTP(pCtx context.Context, wg *sync.WaitGroup, pChErr chan<- error) {
	defer wg.Done()

	go func() {
		err := p.fServiceHTTP.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			pChErr <- err
			return
		}
	}()

	<-pCtx.Done()
}

func (p *sApp) stop() error {
	err := closer.CloseAll([]types.ICloser{
		p.fServiceHTTP,
		p.fServicePPROF,
	})
	if err != nil {
		return utils.MergeErrors(ErrClose, err)
	}
	return nil
}
