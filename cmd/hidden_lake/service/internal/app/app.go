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
	"github.com/number571/go-peer/pkg/storage/database"
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

	fConfig      config.IConfig
	fNode        anonymity.INode
	fConnKeeper  conn_keeper.IConnKeeper
	fServiceHTTP *http.Server
}

func NewApp(
	cfg config.IConfig,
	privKey asymmetric.IPrivKey,
) types.ICommand {
	node := initNode(cfg, privKey)
	return &sApp{
		fConfig:      cfg,
		fNode:        node,
		fConnKeeper:  initConnKeeper(cfg, node),
		fServiceHTTP: initServiceHTTP(config.NewWrapper(cfg), node),
	}
}

func (app *sApp) Run() error {
	app.fMutex.Lock()
	defer app.fMutex.Unlock()

	if app.fIsRun {
		return errors.New("application already running")
	}
	app.fIsRun = true

	res := make(chan error)
	app.fNode.GetWrapperDB().Set(database.NewLevelDB(
		database.NewSettings(&database.SSettings{
			FPath:    pkg_settings.CPathDB,
			FHashing: true,
		}),
	))

	go func() {
		if app.fConfig.GetAddress().GetHTTP() == "" {
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
		app.fNode.HandleFunc(
			pkg_settings.CHeaderHLS,
			handler.HandleServiceTCP(app.fConfig),
		)
		if err := app.fNode.Run(); err != nil {
			res <- err
			return
		}

		// if node in client mode
		// then run endless loop
		tcpAddress := app.fConfig.GetAddress().GetTCP()
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
		app.fNode.GetLogger().PushInfo("HLS is running...")
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

	app.fNode.GetLogger().PushInfo("HLS is shutting down...")
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
