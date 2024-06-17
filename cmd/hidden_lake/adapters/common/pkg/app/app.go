package app

import (
	"context"
	"fmt"
	"sync"

	"github.com/number571/go-peer/cmd/hidden_lake/adapters/common/internal/config"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/state"
	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/go-peer/pkg/utils"

	hla_settings "github.com/number571/go-peer/cmd/hidden_lake/adapters/common/pkg/settings"
	"github.com/number571/go-peer/internal/logger/std"
	internal_types "github.com/number571/go-peer/internal/types"
)

var (
	_ types.IRunner = &sApp{}
)

type sApp struct {
	fState      state.IState
	fStdfLogger logger.ILogger
	fConsumer   types.IRunner
	fProducer   types.IRunner
}

func NewApp(
	pCfg config.IConfig,
	pConsumer types.IRunner,
	pProducer types.IRunner,
) types.IRunner {
	return &sApp{
		fState:      state.NewBoolState(),
		fStdfLogger: std.NewStdLogger(pCfg.GetLogging(), std.GetLogFunc()),
		fConsumer:   pConsumer,
		fProducer:   pProducer,
	}
}

func (p *sApp) Run(pCtx context.Context) error {
	services := []internal_types.IServiceF{
		p.runConsumer,
		p.runProducer,
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
		p.fStdfLogger.PushInfo(fmt.Sprintf( // nolint: perfsprint
			"%s is started",
			hla_settings.CServiceName,
		))
		return nil
	}
}

func (p *sApp) disable(pCancel context.CancelFunc, pWg *sync.WaitGroup) state.IStateF {
	return func() error {
		pCancel()
		pWg.Wait() // wait canceled context

		p.fStdfLogger.PushInfo(fmt.Sprintf( // nolint: perfsprint
			"%s is stopped",
			hla_settings.CServiceName,
		))
		return nil
	}
}

func (p *sApp) runConsumer(pCtx context.Context, wg *sync.WaitGroup, pChErr chan<- error) {
	defer wg.Done()

	if err := p.fConsumer.Run(pCtx); err != nil {
		pChErr <- err
		return
	}
}

func (p *sApp) runProducer(pCtx context.Context, wg *sync.WaitGroup, pChErr chan<- error) {
	defer wg.Done()

	if err := p.fProducer.Run(pCtx); err != nil {
		pChErr <- err
		return
	}
}
