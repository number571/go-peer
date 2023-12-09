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
	"github.com/number571/go-peer/pkg/utils"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/handler"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/interrupt"
	anon_logger "github.com/number571/go-peer/internal/logger/anon"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	std_logger "github.com/number571/go-peer/internal/logger/std"
)

const (
	cInitStart = 3 * time.Second
)

var (
	_ types.IApp = &sApp{}
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
) types.IApp {
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

func (p *sApp) Run() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if p.fIsRun {
		return errors.New("application already running")
	}
	p.fIsRun = true

	if err := p.initDatabase(); err != nil {
		return fmt.Errorf("init database: %w", err)
	}

	p.initServiceHTTP()
	p.initServicePPROF()

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
				p.fAnonLogger,
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
		if err != nil && !errors.Is(err, net.ErrClosed) {
			res <- err
			return
		}
	}()

	select {
	case err := <-res:
		resErr := fmt.Errorf("got run error: %w", err)
		return utils.MergeErrors(resErr, p.Stop())
	case <-time.After(cInitStart):
		p.fStdfLogger.PushInfo(fmt.Sprintf("%s is running...", pkg_settings.CServiceName))
		return nil
	}
}

func (p *sApp) Stop() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if !p.fIsRun {
		return errors.New("application already stopped or not started")
	}
	p.fIsRun = false
	p.fStdfLogger.PushInfo(fmt.Sprintf("%s is shutting down...", pkg_settings.CServiceName))

	err := utils.MergeErrors(
		interrupt.StopAll([]types.IApp{
			p.fNode,
			p.fConnKeeper,
			p.fNode.GetNetworkNode(),
		}),
		interrupt.CloseAll([]types.ICloser{
			p.fServiceHTTP,
			p.fServicePPROF,
			p.fNode.GetWrapperDB(),
		}),
	)
	if err != nil {
		return fmt.Errorf("close/stop all: %w", err)
	}
	return nil
}
