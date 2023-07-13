package app

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/app/state"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/database"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/types"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/settings"
	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/client"
	internal_logger "github.com/number571/go-peer/internal/logger"
	pkg_errors "github.com/number571/go-peer/pkg/errors"
)

const (
	cInitStart = time.Second * 3
)

var (
	_ types.ICommand = &sApp{}
)

type sApp struct {
	fIsRun bool
	fMutex sync.Mutex

	fConfig         config.IConfig
	fStateManager   state.IStateManager
	fLogger         logger.ILogger
	fIntServiceHTTP *http.Server
	fIncServiceHTTP *http.Server
}

func NewApp(
	pCfg config.IConfig,
	pPathTo string,
) types.ICommand {
	stg, err := initCryptoStorage(pCfg, pPathTo)
	if err != nil {
		panic(err)
	}

	state := state.NewStateManager(
		stg,
		database.NewWrapperDB(),
		hls_client.NewClient(
			hls_client.NewBuilder(),
			hls_client.NewRequester(
				fmt.Sprintf("http://%s", pCfg.GetConnection().GetService()),
				&http.Client{Timeout: time.Minute},
			),
		),
		hlt_client.NewClient(
			hlt_client.NewBuilder(),
			hlt_client.NewRequester(
				fmt.Sprintf("http://%s", pCfg.GetConnection().GetTraffic()),
				&http.Client{Timeout: time.Minute},
				message.NewSettings(&message.SSettings{
					FWorkSize:    pCfg.GetWorkSize(),
					FMessageSize: pCfg.GetMessageSize(),
				}),
			),
		),
		pPathTo,
	)

	return &sApp{
		fConfig:       pCfg,
		fStateManager: state,
		fLogger:       internal_logger.StdLogger(pCfg.GetLogging()),
	}
}

func (p *sApp) Run() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if p.fIsRun {
		return pkg_errors.NewError("application already running")
	}
	p.fIsRun = true

	p.initIncomingServiceHTTP()
	p.initInterfaceServiceHTTP()

	res := make(chan error)

	go func() {
		err := p.fIntServiceHTTP.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			res <- err
			return
		}
	}()

	go func() {
		err := p.fIncServiceHTTP.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
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

	// state may be already closed by HLS
	_ = p.fStateManager.CloseState()

	err := types.CloseAll([]types.ICloser{
		p.fIntServiceHTTP,
		p.fIncServiceHTTP,
		p.fStateManager.GetWrapperDB(),
	})
	if err != nil {
		return pkg_errors.WrapError(err, "close/stop all")
	}
	return nil
}
