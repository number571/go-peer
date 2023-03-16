package app

import (
	"errors"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/network/conn_keeper"
	"github.com/number571/go-peer/pkg/types"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/handler"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

const (
	initStart = 3 * time.Second
)

var (
	_ types.ICommand = &sApp{}
)

type sApp struct {
	fIsRun bool
	fMutex sync.Mutex

	fConfig  config.IConfig
	fPrivKey asymmetric.IPrivKey
	fWrapper config.IWrapper

	fNode        anonymity.INode
	fConnKeeper  conn_keeper.IConnKeeper
	fServiceHTTP *http.Server
}

func NewApp(
	cfg config.IConfig,
	privKey asymmetric.IPrivKey,
) types.ICommand {
	return &sApp{
		fConfig:  cfg,
		fPrivKey: privKey,
		fWrapper: config.NewWrapper(cfg),
	}
}

func (app *sApp) Run() error {
	app.fMutex.Lock()
	defer app.fMutex.Unlock()

	if app.fIsRun {
		return errors.New("application already running")
	}
	app.fIsRun = true

	// need reload node for close database in the stop method
	app.fNode = initNode(app.fConfig, app.fPrivKey)
	app.fConnKeeper = initConnKeeper(app.fConfig, app.fNode)
	app.fServiceHTTP = initServiceHTTP(app.fWrapper, app.fNode)

	res := make(chan error)

	go func() {
		if app.fWrapper.GetConfig().GetAddress().GetHTTP() == "" {
			return
		}

		err := app.fServiceHTTP.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			res <- err
			return
		}
	}()

	go func() {
		if err := app.fConnKeeper.Run(); err != nil {
			res <- err
			return
		}
	}()

	go func() {
		cfg := app.fWrapper.GetConfig()

		app.fNode.HandleFunc(
			pkg_settings.CHeaderHLS,
			handler.HandleServiceTCP(cfg),
		)
		if err := app.fNode.Run(); err != nil {
			res <- err
			return
		}

		// if node in client mode
		// then run endless loop
		tcpAddress := cfg.GetAddress().GetTCP()
		if tcpAddress == "" {
			select {}
		}

		// run node in server mode
		err := app.fNode.GetNetworkNode().Run()
		if err != nil && !errors.Is(err, net.ErrClosed) {
			res <- err
			return
		}
	}()

	select {
	case err := <-res:
		app.Stop()
		return err
	case <-time.After(initStart):
		app.fNode.GetLogger().PushInfo("service is running...")
		return nil
	}
}

func (app *sApp) Stop() error {
	app.fMutex.Lock()
	defer app.fMutex.Unlock()

	if !app.fIsRun {
		return errors.New("application already stopped or not started")
	}
	app.fIsRun = false

	app.fNode.GetLogger().PushInfo("service is shutting down...")
	app.fNode.HandleFunc(pkg_settings.CHeaderHLS, nil)

	lastErr := types.StopAll([]types.ICommand{
		app.fNode,
		app.fConnKeeper,
		app.fNode.GetNetworkNode(),
	})

	err := types.CloseAll([]types.ICloser{
		app.fServiceHTTP,
		app.fNode.GetWrapperDB(),
	})
	if err != nil {
		lastErr = err
	}

	return lastErr
}
