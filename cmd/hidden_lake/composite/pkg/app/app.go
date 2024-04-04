package app

import (
	"context"
	"sync"

	"github.com/number571/go-peer/cmd/hidden_lake/composite/internal/config"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/state"
	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/go-peer/pkg/utils"

	std_logger "github.com/number571/go-peer/internal/logger/std"
	internal_types "github.com/number571/go-peer/internal/types"

	hlc_settings "github.com/number571/go-peer/cmd/hidden_lake/composite/pkg/settings"
)

var (
	_ types.IRunner = &sApp{}
)

type sApp struct {
	fState   state.IState
	fRunners []types.IRunner

	fStdfLogger logger.ILogger
}

func NewApp(
	pCfg config.IConfig,
	pRunners []types.IRunner,
) types.IRunner {
	stdfLogger := std_logger.NewStdLogger(pCfg.GetLogging(), std_logger.GetLogFunc())

	return &sApp{
		fState:      state.NewBoolState(),
		fRunners:    pRunners,
		fStdfLogger: stdfLogger,
	}
}

func (p *sApp) Run(pCtx context.Context) error {
	services := runnersToServices(p.fRunners)

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
		p.fStdfLogger.PushInfo(hlc_settings.CServiceName + " is running...")
		return nil
	}
}

func (p *sApp) disable(pCancel context.CancelFunc, pWg *sync.WaitGroup) state.IStateF {
	return func() error {
		pCancel()
		pWg.Wait() // wait canceled context

		p.fStdfLogger.PushInfo(hlc_settings.CServiceName + " is shutting down...")
		return nil
	}
}

func runnersToServices(pRunners []types.IRunner) []internal_types.IServiceF {
	services := make([]internal_types.IServiceF, 0, len(pRunners))
	for _, r := range pRunners {
		r := r
		services = append(services, runnerToService(r))
	}
	return services
}

func runnerToService(pRunner types.IRunner) internal_types.IServiceF {
	return func(pCtx context.Context, pWg *sync.WaitGroup, pChErr chan<- error) {
		defer pWg.Done()
		if err := pRunner.Run(pCtx); err != nil {
			pChErr <- err
			return
		}
	}
}
