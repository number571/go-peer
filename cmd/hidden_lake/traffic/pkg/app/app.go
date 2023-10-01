package app

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn_keeper"
	"github.com/number571/go-peer/pkg/types"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	anon_logger "github.com/number571/go-peer/internal/logger/anon"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	std_logger "github.com/number571/go-peer/internal/logger/std"
	pkg_errors "github.com/number571/go-peer/pkg/errors"
)

const (
	cInitStart = 3 * time.Second
)

var (
	_ types.ICommand = &sApp{}
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
) types.ICommand {
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

func (p *sApp) Run() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if p.fIsRun {
		return pkg_errors.NewError("application already running")
	}
	p.fIsRun = true

	p.initServiceHTTP()
	p.initServicePPROF()

	if err := p.initDatabase(); err != nil {
		return err
	}

	res := make(chan error)

	go func() {
		if err := p.fConnKeeper.Run(); err != nil {
			res <- err
			return
		}
	}()

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

	go func() {
		// if node in client mode
		tcpAddress := p.fConfig.GetAddress().GetTCP()
		if tcpAddress == "" {
			return
		}

		if err := p.fNode.Run(); err != nil {
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

	err := pkg_errors.AppendError(
		types.StopAll([]types.ICommand{
			p.fConnKeeper,
			p.fConnKeeper.GetNetworkNode(),
		}),
		types.CloseAll([]types.ICloser{
			p.fServiceHTTP,
			p.fServicePPROF,
			p.fWrapperDB,
		}),
	)
	if err != nil {
		return pkg_errors.WrapError(err, "close/stop all")
	}
	return nil
}
