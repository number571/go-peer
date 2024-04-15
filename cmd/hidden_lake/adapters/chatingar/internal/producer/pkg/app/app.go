package app

import (
	"context"

	"github.com/number571/go-peer/cmd/hidden_lake/adapters"
	"github.com/number571/go-peer/cmd/hidden_lake/adapters/chatingar/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/adapters/chatingar/internal/producer/internal/adapted"
	"github.com/number571/go-peer/internal/logger/std"
	"github.com/number571/go-peer/pkg/logger"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/state"
	"github.com/number571/go-peer/pkg/types"
)

var (
	_ types.IRunner = &sApp{}
)

type sApp struct {
	fState        state.IState
	fStdfLogger   logger.ILogger
	fSettings     net_message.ISettings
	fPostID       string
	fIncomingAddr string
}

func NewApp(pCfg config.IConfig) types.IRunner {
	return &sApp{
		fState:        state.NewBoolState(),
		fStdfLogger:   std.NewStdLogger(pCfg.GetLogging(), std.GetLogFunc()),
		fPostID:       pCfg.GetConnection().GetPostID(),
		fIncomingAddr: pCfg.GetAddress(),
		fSettings:     pCfg.GetSettings(),
	}
}

func (p *sApp) Run(pCtx context.Context) error {
	if err := p.fState.Enable(nil); err != nil {
		return err
	}
	defer func() { _ = p.fState.Disable(nil) }()

	return adapters.ProduceProcessor(
		pCtx,
		adapted.NewAdaptedProducer(p.fPostID),
		p.fStdfLogger,
		p.fSettings,
		p.fIncomingAddr,
	)
}
