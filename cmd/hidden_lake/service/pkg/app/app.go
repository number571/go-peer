package app

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/network/conn_keeper"
	"github.com/number571/go-peer/pkg/types"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/handler"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
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
	fWrapper    config.IWrapper
	fNode       anonymity.INode
	fLogger     logger.ILogger
	fConnKeeper conn_keeper.IConnKeeper

	fServiceHTTP  *http.Server
	fServicePPROF *http.Server
}

func NewApp(
	pCfg config.IConfig,
	pPrivKey asymmetric.IPrivKey,
	pPathTo string,
) types.ICommand {
	logger := internal_logger.StdLogger(pCfg.GetLogging())
	node := initNode(pCfg, pPrivKey, logger)
	return &sApp{
		fPathTo:     pPathTo,
		fWrapper:    config.NewWrapper(pCfg),
		fNode:       node,
		fLogger:     logger,
		fConnKeeper: initConnKeeper(pCfg, node),
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
		return pkg_errors.WrapError(err, "init database")
	}

	res := make(chan error)

	go func() {
		if p.fWrapper.GetConfig().GetAddress().GetPPROF() == "" {
			return
		}

		err := p.fServicePPROF.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			res <- err
			return
		}
	}()

	go func() {
		if p.fWrapper.GetConfig().GetAddress().GetHTTP() == "" {
			return
		}

		err := p.fServiceHTTP.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			res <- err
			return
		}
	}()

	go func() {
		if err := p.fConnKeeper.Run(); err != nil {
			res <- err
			return
		}
	}()

	go func() {
		p.fNode.HandleFunc(
			pkg_settings.CServiceMask,
			handler.HandleServiceTCP(
				p.fWrapper.GetConfig(),
				p.fLogger,
			),
		)
		if err := p.fNode.Run(); err != nil {
			res <- err
			return
		}

		// if node in client mode
		tcpAddress := p.fWrapper.GetConfig().GetAddress().GetTCP()
		if tcpAddress == "" {
			return
		}

		// run node in server mode
		err := p.fNode.GetNetworkNode().Run()
		if err != nil && !pkg_errors.HasError(err, net.ErrClosed) {
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

	p.fNode.HandleFunc(pkg_settings.CServiceMask, nil)
	err := pkg_errors.AppendError(
		types.StopAll([]types.ICommand{
			p.fNode,
			p.fConnKeeper,
			p.fNode.GetNetworkNode(),
		}),
		types.CloseAll([]types.ICloser{
			p.fServiceHTTP,
			p.fServicePPROF,
			p.fNode.GetWrapperDB(),
		}),
	)
	if err != nil {
		return pkg_errors.WrapError(err, "close/stop all")
	}
	return nil
}
