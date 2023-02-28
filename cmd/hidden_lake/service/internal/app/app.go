package app

import (
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
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
	fWrapper     config.IWrapper
	fNode        anonymity.INode
	fConnKeeper  conn_keeper.IConnKeeper
	fServiceHTTP *http.Server
}

func NewApp(
	cfg config.IConfig,
	node anonymity.INode,
) types.ICommand {
	wrapper := config.NewWrapper(cfg)
	return &sApp{
		fWrapper:     wrapper,
		fNode:        node,
		fConnKeeper:  initConnKeeper(cfg, node),
		fServiceHTTP: initServiceHTTP(wrapper, node),
	}
}

func (app *sApp) Run() error {
	res := make(chan error)

	go func() {
		if app.fWrapper.Config().Address().HTTP() == "" {
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
		cfg := app.fWrapper.Config()

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
		tcpAddress := cfg.Address().TCP()
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
		return nil
	}
}

func (app *sApp) Stop() error {
	app.fNode.HandleFunc(pkg_settings.CHeaderHLS, nil)

	lastErr := types.StopAllCommands([]types.ICommand{
		app.fNode,
		app.fConnKeeper,
		app.fNode.GetNetworkNode(),
	})

	err := types.CloseAll([]types.ICloser{
		app.fServiceHTTP,
		app.fNode.GetKeyValueDB(),
	})
	if err != nil {
		lastErr = err
	}

	return lastErr
}
