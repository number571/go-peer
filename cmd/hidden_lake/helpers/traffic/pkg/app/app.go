package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/internal/database"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/connkeeper"
	"github.com/number571/go-peer/pkg/state"
	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/go-peer/pkg/utils"

	pkg_config "github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/pkg/config"
	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/pkg/settings"
	"github.com/number571/go-peer/internal/closer"
	anon_logger "github.com/number571/go-peer/internal/logger/anon"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	std_logger "github.com/number571/go-peer/internal/logger/std"
	internal_types "github.com/number571/go-peer/internal/types"
)

var (
	_ types.IRunner = &sApp{}
)

type sApp struct {
	fState state.IState

	fPathTo     string
	fConfig     config.IConfig
	fNode       network.INode
	fDatabase   database.IDatabase
	fConnKeeper connkeeper.IConnKeeper

	fAnonLogger logger.ILogger
	fHTTPLogger logger.ILogger
	fStdfLogger logger.ILogger

	fServiceHTTP  *http.Server
	fServicePPROF *http.Server
}

func NewApp(
	pCfg config.IConfig,
	pPathTo string,
) types.IRunner {
	var (
		anonLogger = std_logger.NewStdLogger(pCfg.GetLogging(), anon_logger.GetLogFunc())
		httpLogger = std_logger.NewStdLogger(pCfg.GetLogging(), http_logger.GetLogFunc())
		stdfLogger = std_logger.NewStdLogger(pCfg.GetLogging(), std_logger.GetLogFunc())
	)

	return &sApp{
		fState:      state.NewBoolState(),
		fConfig:     pCfg,
		fPathTo:     pPathTo,
		fAnonLogger: anonLogger,
		fHTTPLogger: httpLogger,
		fStdfLogger: stdfLogger,
	}
}

func (p *sApp) Run(pCtx context.Context) error {
	services := []internal_types.IServiceF{
		p.runListenerPPROF,
		p.runListenerHTTP,
		p.runListenerNode,
		p.runConnKeeper,
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

func (p *sApp) enable(pCtx context.Context) state.IStateF {
	return func() error {
		if err := p.initDatabase(); err != nil {
			return utils.MergeErrors(ErrInitDB, err)
		}

		p.initNetworkNode(p.fDatabase)
		p.initConnKeeper(p.fNode)

		p.initServicePPROF()
		p.initServiceHTTP(pCtx)

		p.fStdfLogger.PushInfo(fmt.Sprintf(
			"%s is started; %s",
			hlt_settings.CServiceName,
			encoding.SerializeJSON(pkg_config.GetConfigSettings(p.fConfig)),
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
			hlt_settings.CServiceName,
		))
		return p.stop()
	}
}
func (p *sApp) stop() error {
	err := closer.CloseAll([]types.ICloser{
		p.fServiceHTTP,
		p.fServicePPROF,
		p.fDatabase,
		p.fNode,
	})
	if err != nil {
		return utils.MergeErrors(ErrClose, err)
	}
	return nil
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

	if p.fConfig.GetAddress().GetHTTP() == "" {
		return
	}

	go func() {
		err := p.fServiceHTTP.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			pChErr <- err
			return
		}
	}()

	<-pCtx.Done()
}

func (p *sApp) runListenerNode(pCtx context.Context, wg *sync.WaitGroup, pChErr chan<- error) {
	defer wg.Done()

	// if node in client mode
	tcpAddress := p.fConfig.GetAddress().GetTCP()
	if tcpAddress == "" {
		return
	}

	go func() {
		// run node in server mode
		err := p.fNode.Listen(pCtx)
		if err != nil && !errors.Is(err, net.ErrClosed) {
			pChErr <- err
			return
		}
	}()

	<-pCtx.Done()
}

func (p *sApp) runConnKeeper(pCtx context.Context, wg *sync.WaitGroup, pChErr chan<- error) {
	defer wg.Done()

	if err := p.fConnKeeper.Run(pCtx); err != nil {
		pChErr <- err
		return
	}
}
