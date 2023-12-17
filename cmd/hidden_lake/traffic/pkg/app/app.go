package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn_keeper"
	"github.com/number571/go-peer/pkg/types"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	"github.com/number571/go-peer/internal/interrupt"
	anon_logger "github.com/number571/go-peer/internal/logger/anon"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	std_logger "github.com/number571/go-peer/internal/logger/std"
	"github.com/number571/go-peer/pkg/utils"
)

var (
	_ types.IRunner = &sApp{}
)

type sApp struct {
	fIsRun bool
	fMutex sync.Mutex

	fPathTo     string
	fConfig     config.IConfig
	fWrapperDB  database.IWrapperDB
	fNode       network.INode
	fConnKeeper conn_keeper.IConnKeeper

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
	anonLogger := std_logger.NewStdLogger(
		pCfg.GetLogging(),
		anon_logger.GetLogFunc(),
	)

	httpLogger := std_logger.NewStdLogger(
		pCfg.GetLogging(),
		http_logger.GetLogFunc(),
	)

	stdfLogger := std_logger.NewStdLogger(
		pCfg.GetLogging(),
		std_logger.GetLogFunc(),
	)

	wDB := database.NewWrapperDB()
	node := initNode(pCfg, wDB, anonLogger)

	return &sApp{
		fConfig:     pCfg,
		fWrapperDB:  wDB,
		fNode:       node,
		fConnKeeper: initConnKeeper(pCfg, node),
		fPathTo:     pPathTo,
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
		if err := p.fConnKeeper.Run(pCtx); err != nil {
			chErr <- err
			return
		}
	}()

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

	go func() {
		// if node in client mode
		tcpAddress := p.fConfig.GetAddress().GetTCP()
		if tcpAddress == "" {
			return
		}

		if err := p.fNode.Listen(); err != nil {
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
			p.fWrapperDB,
			p.fConnKeeper.GetNetworkNode(),
		}),
	)
	if err != nil {
		return fmt.Errorf("close/stop all: %w", err)
	}

	return nil
}
