package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/network/conn_keeper"
	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/go-peer/pkg/utils"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/interrupt"
	anon_logger "github.com/number571/go-peer/internal/logger/anon"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	std_logger "github.com/number571/go-peer/internal/logger/std"
)

var (
	_ types.IRunner = &sApp{}
)

type sApp struct {
	fIsRun bool
	fMutex sync.Mutex

	fPathTo     string
	fWrapper    config.IWrapper
	fNode       anonymity.INode
	fConnKeeper conn_keeper.IConnKeeper
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
) types.IRunner {
	logging := pCfg.GetLogging()

	var (
		anonLogger = std_logger.NewStdLogger(logging, anon_logger.GetLogFunc())
		httpLogger = std_logger.NewStdLogger(logging, http_logger.GetLogFunc())
		stdfLogger = std_logger.NewStdLogger(logging, std_logger.GetLogFunc())
	)

	node := initNode(pCfg, pPrivKey, anonLogger)
	return &sApp{
		fPathTo:     pPathTo,
		fWrapper:    config.NewWrapper(pCfg),
		fNode:       node,
		fConnKeeper: initConnKeeper(pCfg, node),
		fPrivKey:    pPrivKey,
		fAnonLogger: anonLogger,
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

	p.initServiceHTTP()
	p.initServicePPROF()

	chErr := make(chan error)

	go func() {
		if p.fWrapper.GetConfig().GetAddress().GetPPROF() == "" {
			return
		}

		err := p.fServicePPROF.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			chErr <- err
			return
		}
	}()

	go func() {
		if p.fWrapper.GetConfig().GetAddress().GetHTTP() == "" {
			return
		}

		err := p.fServiceHTTP.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			chErr <- err
			return
		}
	}()

	go func() {
		// if node in client mode
		tcpAddress := p.fWrapper.GetConfig().GetAddress().GetTCP()
		if tcpAddress == "" {
			return
		}

		// run node in server mode
		err := p.fNode.GetNetworkNode().Listen()
		if err != nil && !errors.Is(err, net.ErrClosed) {
			chErr <- err
			return
		}
	}()

	go func() {
		if err := p.fConnKeeper.Run(pCtx); err != nil {
			chErr <- err
			return
		}
	}()

	go func() {
		if err := p.fNode.Run(pCtx); err != nil {
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

	err := utils.MergeErrors(
		interrupt.CloseAll([]types.ICloser{
			p.fServiceHTTP,
			p.fServicePPROF,
			p.fNode.GetWrapperDB(),
			p.fNode.GetNetworkNode(),
		}),
	)
	if err != nil {
		return fmt.Errorf("close/stop all: %w", err)
	}

	return nil
}
