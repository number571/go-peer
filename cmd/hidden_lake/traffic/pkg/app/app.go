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
	internal_logger "github.com/number571/go-peer/internal/logger/std"
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
	fLogger     logger.ILogger
	fNode       network.INode
	fConnKeeper conn_keeper.IConnKeeper

	fServiceHTTP  *http.Server
	fServicePPROF *http.Server
}

func NewApp(
	pCfg config.IConfig,
	pPathTo string,
) types.ICommand {
	wDB := database.NewWrapperDB()
	logger := internal_logger.StdLogger(pCfg.GetLogging())
	node := initNode(pCfg, wDB, logger)
	connKeeper := initConnKeeper(pCfg, node)
	return &sApp{
		fConfig:     pCfg,
		fWrapperDB:  wDB,
		fLogger:     logger,
		fNode:       node,
		fConnKeeper: connKeeper,
		fPathTo:     pPathTo,
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
		p.fLogger.PushInfo(fmt.Sprintf("%s is running...", pkg_settings.CServiceName))
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
	p.fLogger.PushInfo(fmt.Sprintf("%s is shutting down...", pkg_settings.CServiceName))

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
