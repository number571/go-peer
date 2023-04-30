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
	"github.com/number571/go-peer/pkg/network/conn_keeper"
	"github.com/number571/go-peer/pkg/types"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
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
	fWrapperDB   database.IWrapperDB
	fLogger      logger.ILogger
	fConnKeeper  conn_keeper.IConnKeeper
	fServiceHTTP *http.Server
}

func NewApp(
	pCfg config.IConfig,
	pPathTo string,
) types.ICommand {
	wDB := database.NewWrapperDB()
	logger := internal_logger.StdLogger(pCfg.GetLogging())
	connKeeper := initConnKeeper(pCfg, wDB, logger)
	return &sApp{
		fConfig:     pCfg,
		fWrapperDB:  wDB,
		fLogger:     logger,
		fConnKeeper: connKeeper,
		fPathTo:     pPathTo,
	}
}

func (p *sApp) Run() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if p.fIsRun {
		return errors.New("application already running")
	}
	p.fIsRun = true

	p.initServiceHTTP()
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
		if p.fConfig.GetAddress() == "" {
			return
		}

		err := p.fServiceHTTP.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
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

	lastErr := types.StopAll([]types.ICommand{
		p.fConnKeeper,
		p.fConnKeeper.GetNetworkNode(),
	})

	err := types.CloseAll([]types.ICloser{
		p.fServiceHTTP,
		p.fWrapperDB,
	})
	if err != nil {
		lastErr = err
	}

	return lastErr
}
