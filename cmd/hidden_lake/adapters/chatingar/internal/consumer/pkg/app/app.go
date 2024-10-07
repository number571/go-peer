package app

import (
	"context"
	"net/http"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/adapters"
	"github.com/number571/go-peer/cmd/hidden_lake/adapters/chatingar/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/adapters/chatingar/internal/consumer/internal/adapted"
	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/pkg/client"
	"github.com/number571/go-peer/internal/logger/std"
	"github.com/number571/go-peer/pkg/logger"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/state"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/go-peer/pkg/types"
)

var (
	_ types.IRunner = &sApp{}
)

type sApp struct {
	fState      state.IState
	fStdfLogger logger.ILogger
	fSettings   net_message.ISettings
	fHltAddr    string
	fPostID     string
	fWaitTime   time.Duration
}

func NewApp(pCfg config.IConfig) types.IRunner {
	return &sApp{
		fState:      state.NewBoolState(),
		fStdfLogger: std.NewStdLogger(pCfg.GetLogging(), std.GetLogFunc()),
		fHltAddr:    pCfg.GetConnection().GetHLTHost(),
		fPostID:     pCfg.GetConnection().GetPostID(),
		fSettings:   pCfg.GetSettings(),
		fWaitTime:   time.Duration(pCfg.GetSettings().GetWaitTimeMS()) * time.Millisecond,
	}
}

func (p *sApp) Run(pCtx context.Context) error {
	if err := p.fState.Enable(nil); err != nil {
		return err
	}
	defer func() { _ = p.fState.Disable(nil) }()

	return adapters.ConsumeProcessor(
		pCtx,
		adapted.NewAdaptedConsumer(
			p.fPostID,
			p.fSettings,
			cache.NewLRUCache(1<<10),
		),
		p.fStdfLogger,
		hlt_client.NewClient(
			hlt_client.NewBuilder(),
			hlt_client.NewRequester(
				"http://"+p.fHltAddr,
				&http.Client{Timeout: time.Minute},
				p.fSettings,
			),
		),
		p.fWaitTime,
	)
}
