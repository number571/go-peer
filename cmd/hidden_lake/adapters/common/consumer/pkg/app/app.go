package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/adapters"
	"github.com/number571/go-peer/cmd/hidden_lake/adapters/common/consumer/internal/adapted"
	"github.com/number571/go-peer/cmd/hidden_lake/adapters/common/consumer/internal/config"
	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/pkg/client"
	"github.com/number571/go-peer/internal/logger/std"
	"github.com/number571/go-peer/pkg/database"
	"github.com/number571/go-peer/pkg/logger"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/state"
	"github.com/number571/go-peer/pkg/types"
)

const (
	cDBPath = "common_consumer.db"
)

var (
	_ types.IRunner = &sApp{}
)

type sApp struct {
	fState       state.IState
	fStdfLogger  logger.ILogger
	fSettings    net_message.ISettings
	fHltAddr     string
	fServiceAddr string
	fWaitTime    time.Duration
}

func NewApp(pCfg config.IConfig) types.IRunner {
	return &sApp{
		fState:       state.NewBoolState(),
		fStdfLogger:  std.NewStdLogger(pCfg.GetLogging(), std.GetLogFunc()),
		fHltAddr:     pCfg.GetConnection().GetHLTHost(),
		fServiceAddr: pCfg.GetConnection().GetSrvHost(),
		fSettings:    pCfg.GetSettings(),
		fWaitTime:    time.Duration(pCfg.GetSettings().GetWaitTimeMS()) * time.Millisecond,
	}
}

func (p *sApp) Run(pCtx context.Context) error {
	if err := p.fState.Enable(p.enable()); err != nil {
		return err
	}
	defer func() { _ = p.fState.Disable(p.disable()) }()

	kvDB, err := database.NewKVDatabase(
		database.NewSettings(&database.SSettings{
			FPath: cDBPath,
		}),
	)
	if err != nil {
		return err
	}
	defer kvDB.Close()

	return adapters.ConsumeProcessor(
		pCtx,
		adapted.NewAdaptedConsumer(p.fServiceAddr, p.fSettings, kvDB),
		p.fStdfLogger,
		hlt_client.NewClient(
			hlt_client.NewBuilder(),
			hlt_client.NewRequester(
				fmt.Sprintf("http://%s", p.fHltAddr),
				&http.Client{Timeout: time.Minute},
				p.fSettings,
			),
		),
		p.fWaitTime,
	)
}

func (p *sApp) enable() state.IStateF {
	return func() error {
		p.fStdfLogger.PushInfo("common/consumer is running...")
		return nil
	}
}

func (p *sApp) disable() state.IStateF {
	return func() error {
		p.fStdfLogger.PushInfo("common/consumer is shutting down...")
		return nil
	}
}
