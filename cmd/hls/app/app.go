package app

import (
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/number571/go-peer/cmd/hls/config"
	"github.com/number571/go-peer/modules"
	"github.com/number571/go-peer/modules/closer"
	"github.com/number571/go-peer/modules/network/anonymity"
	"github.com/number571/go-peer/modules/network/conn_keeper"

	"github.com/number571/go-peer/cmd/hls/handler"
	hls_settings "github.com/number571/go-peer/cmd/hls/settings"
)

const (
	initStart = 3 * time.Second
)

var (
	_ IApp = &sApp{}
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
) IApp {
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
		httpAddress := app.fWrapper.Config().Address().HTTP()
		if httpAddress == "" {
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

		app.fNode.Handle(
			hls_settings.CHeaderHLS,
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
		err := app.fNode.Network().Listen(tcpAddress)
		if err != nil && !errors.Is(err, net.ErrClosed) {
			res <- err
			return
		}
	}()

	select {
	case err := <-res:
		app.Close()
		return err
	case <-time.After(initStart):
		return nil
	}
}

func (app *sApp) Close() error {
	app.fNode.Handle(hls_settings.CHeaderHLS, nil)
	return closer.CloseAll([]modules.ICloser{
		app.fNode,
		app.fConnKeeper,
		app.fServiceHTTP,
		app.fNode.Network(),
		app.fNode.KeyValueDB(),
	})
}
