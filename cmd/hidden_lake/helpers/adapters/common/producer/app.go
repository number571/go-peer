package main

import (
	"context"
	"errors"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/adapters"
	"github.com/number571/go-peer/cmd/hidden_lake/helpers/adapters/common"
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
	fServiceAddr  string
	fIncomingAddr string
}

func initApp(pArgs []string) (types.IRunner, error) {
	if len(pArgs) != 3 {
		return nil, errors.New("./producer [incoming-addr] [service-addr] [logger]")
	}

	incomingAddr := pArgs[0]
	serviceAddr := pArgs[1]
	logEnabled := pArgs[2]

	logList := []string{}
	if logEnabled == "true" {
		logList = []string{"info", "warn", "erro"}
	}

	logging, err := std.LoadLogging(logList)
	if err != nil {
		return nil, err
	}

	return newApp(logging, serviceAddr, incomingAddr), nil
}

func newApp(
	pLogging std.ILogging,
	pServiceAddr string,
	pIncomingAddr string,
) types.IRunner {
	return &sApp{
		fState:        state.NewBoolState(),
		fStdfLogger:   std.NewStdLogger(pLogging, std.GetLogFunc()),
		fSettings:     common.GetMessageSettings(),
		fServiceAddr:  pServiceAddr,
		fIncomingAddr: pIncomingAddr,
	}
}

func (p *sApp) Run(pCtx context.Context) error {
	if err := p.fState.Enable(p.enable()); err != nil {
		return err
	}
	defer func() { _ = p.fState.Disable(p.disable()) }()

	return adapters.ProduceProcessor(
		pCtx,
		newAdaptedProducer(p.fServiceAddr),
		p.fStdfLogger,
		p.fSettings,
		p.fIncomingAddr,
	)
}

func (p *sApp) enable() state.IStateF {
	return func() error {
		p.fStdfLogger.PushInfo("common/producer is running...")
		return nil
	}
}

func (p *sApp) disable() state.IStateF {
	return func() error {
		p.fStdfLogger.PushInfo("common/producer is shutting down...")
		return nil
	}
}
