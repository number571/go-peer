package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	"github.com/number571/go-peer/internal/closer"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/network/connkeeper"
	"github.com/number571/go-peer/pkg/state"
	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/go-peer/pkg/utils"

	pkg_config "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/config"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	anon_logger "github.com/number571/go-peer/internal/logger/anon"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	std_logger "github.com/number571/go-peer/internal/logger/std"
	internal_types "github.com/number571/go-peer/internal/types"
)

var (
	_ types.IRunner = &sApp{}
)

type sApp struct {
	fState    state.IState
	fPathTo   string
	fParallel uint64

	fCfgW       config.IWrapper
	fNode       anonymity.INode
	fConnKeeper connkeeper.IConnKeeper
	fPrivKey    asymmetric.IPrivKey

	fAnonLogger logger.ILogger
	fHTTPLogger logger.ILogger
	fStdfLogger logger.ILogger

	fServiceHTTP  *http.Server
	fServicePPROF *http.Server
}

func NewApp(
	pCfg config.IConfig,
	pPrivKey asymmetric.IPrivKey,
	pPathTo string,
	pParallel uint64,
) types.IRunner {
	logging := pCfg.GetLogging()

	var (
		anonLogger = std_logger.NewStdLogger(logging, anon_logger.GetLogFunc())
		httpLogger = std_logger.NewStdLogger(logging, http_logger.GetLogFunc())
		stdfLogger = std_logger.NewStdLogger(logging, std_logger.GetLogFunc())
	)

	return &sApp{
		fState:      state.NewBoolState(),
		fPathTo:     pPathTo,
		fParallel:   pParallel,
		fCfgW:       config.NewWrapper(pCfg),
		fPrivKey:    pPrivKey,
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
		p.runNode,
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
		if err := p.initAnonNode(); err != nil {
			return utils.MergeErrors(ErrCreateAnonNode, err)
		}

		p.initConnKeeper(
			p.fNode.GetNetworkNode(),
		)

		p.initServicePPROF()
		p.initServiceHTTP(pCtx)

		p.fStdfLogger.PushInfo(fmt.Sprintf(
			"%s is started; %s",
			hls_settings.CServiceName,
			encoding.SerializeJSON(pkg_config.GetConfigSettings(p.fCfgW.GetConfig(), p.fNode)),
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
			hls_settings.CServiceName,
		))
		return p.stop()
	}
}

func (p *sApp) stop() error {
	err := closer.CloseAll([]types.ICloser{
		p.fServiceHTTP,
		p.fServicePPROF,
		p.fNode.GetKVDatabase(),
		p.fNode.GetNetworkNode(),
	})
	if err != nil {
		return utils.MergeErrors(ErrClose, err)
	}
	return nil
}

func (p *sApp) runListenerPPROF(pCtx context.Context, wg *sync.WaitGroup, pChErr chan<- error) {
	defer wg.Done()

	if p.fCfgW.GetConfig().GetAddress().GetPPROF() == "" {
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

	if p.fCfgW.GetConfig().GetAddress().GetHTTP() == "" {
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
	tcpAddress := p.fCfgW.GetConfig().GetAddress().GetTCP()
	if tcpAddress == "" {
		return
	}

	go func() {
		// run node in server mode
		err := p.fNode.GetNetworkNode().Listen(pCtx)
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

func (p *sApp) runNode(pCtx context.Context, wg *sync.WaitGroup, pChErr chan<- error) {
	defer wg.Done()

	if err := p.fNode.Run(pCtx); err != nil {
		pChErr <- err
		return
	}
}
