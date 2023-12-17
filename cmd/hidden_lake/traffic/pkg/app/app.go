package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn_keeper"
	"github.com/number571/go-peer/pkg/state"
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
	fState state.IState

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
		fState:      state.NewBoolState(),
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
	enableFunc := func() error {
		if err := p.initDatabase(); err != nil {
			return fmt.Errorf("init database: %w", err)
		}
		p.fStdfLogger.PushInfo(fmt.Sprintf("%s is running...", pkg_settings.CServiceName))
		return nil
	}
	if err := p.fState.Enable(enableFunc); err != nil {
		return fmt.Errorf("application running error: %w", err)
	}

	defer func() {
		disableFunc := func() error {
			p.fStdfLogger.PushInfo(fmt.Sprintf("%s is shutting down...", pkg_settings.CServiceName))
			return p.stop()
		}
		_ = p.fState.Disable(disableFunc)
	}()

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
		return nil
	case err := <-chErr:
		return fmt.Errorf("got run error: %w", err)
	}
}

func (p *sApp) stop() error {
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
