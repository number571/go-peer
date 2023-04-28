package app

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/config"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/network/conn_keeper"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/go-peer/pkg/types"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/handler"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	internal_logger "github.com/number571/go-peer/internal/logger"
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

	fPathTo      string
	fConfig      config.IConfig
	fNode        anonymity.INode
	fLogger      logger.ILogger
	fConnKeeper  conn_keeper.IConnKeeper
	fServiceHTTP *http.Server
}

func NewApp(
	pCfg config.IConfig,
	pPrivKey asymmetric.IPrivKey,
	pPathTo string,
) types.ICommand {
	logger := internal_logger.StdLogger(pCfg.GetLogging())
	node := initNode(pCfg, pPrivKey, logger)
	return &sApp{
		fConfig:      pCfg,
		fNode:        node,
		fLogger:      logger,
		fPathTo:      pPathTo,
		fConnKeeper:  initConnKeeper(pCfg, node),
		fServiceHTTP: initServiceHTTP(config.NewWrapper(pCfg), node),
	}
}

func (p *sApp) Run() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if p.fIsRun {
		return errors.New("application already running")
	}
	p.fIsRun = true

	db, err := database.NewSQLiteDB(
		database.NewSettings(&database.SSettings{
			FPath:    fmt.Sprintf("%s/%s", p.fPathTo, pkg_settings.CPathDB),
			FHashing: true,
		}),
	)
	if err != nil {
		return err
	}

	res := make(chan error)
	p.fNode.GetWrapperDB().Set(db)

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
		if err := p.fConnKeeper.Run(); err != nil {
			res <- err
			return
		}
	}()

	go func() {
		p.fNode.HandleFunc(
			pkg_settings.CHeaderHLS,
			handler.HandleServiceTCP(p.fConfig),
		)
		if err := p.fNode.Run(); err != nil {
			res <- err
			return
		}

		// if node in client mode
		// then run endless loop
		tcpAddress := p.fConfig.GetAddress().GetTCP()
		if tcpAddress == "" {
			select {}
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
		p.Stop()
		return err
	case <-time.After(cInitStart):
		p.fLogger.PushInfo(fmt.Sprintf("%s is running...", pkg_settings.CServiceName))
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
	p.fLogger.PushInfo(fmt.Sprintf("%s is shutting down...", pkg_settings.CServiceName))

	p.fNode.HandleFunc(pkg_settings.CHeaderHLS, nil)
	lastErr := types.StopAll([]types.ICommand{
		p.fNode,
		p.fConnKeeper,
		p.fNode.GetNetworkNode(),
	})

	err := types.CloseAll([]types.ICloser{
		p.fServiceHTTP,
		p.fNode.GetWrapperDB(),
	})
	if err != nil {
		lastErr = err
	}

	return lastErr
}
