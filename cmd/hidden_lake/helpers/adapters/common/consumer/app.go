package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/adapters"
	"github.com/number571/go-peer/cmd/hidden_lake/helpers/adapters/common"
	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/pkg/client"
	"github.com/number571/go-peer/internal/logger/std"
	"github.com/number571/go-peer/pkg/logger"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/state"
	"github.com/number571/go-peer/pkg/storage"
	"github.com/number571/go-peer/pkg/storage/database"
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
}

func initApp(pArgs []string) (types.IRunner, error) {
	if len(pArgs) != 3 {
		return nil, errors.New("./consumer [service-addr] [hlt-addr] [logger]")
	}

	serviceAddr := pArgs[0]
	hltAddr := pArgs[1]
	logEnabled := pArgs[2]

	logList := []string{}
	if logEnabled == "true" {
		logList = []string{"info", "warn", "erro"}
	}

	logging, err := std.LoadLogging(logList)
	if err != nil {
		return nil, err
	}

	return newApp(logging, serviceAddr, hltAddr), nil
}

func newApp(
	pLogging std.ILogging,
	pServiceAddr string,
	pHltAddr string,
) types.IRunner {
	return &sApp{
		fState:       state.NewBoolState(),
		fStdfLogger:  std.NewStdLogger(pLogging, std.GetLogFunc()),
		fHltAddr:     pHltAddr,
		fSettings:    common.GetMessageSettings(),
		fServiceAddr: pServiceAddr,
	}
}

func (p *sApp) Run(pCtx context.Context) error {
	if err := p.fState.Enable(p.enable()); err != nil {
		return err
	}
	defer func() { _ = p.fState.Disable(p.disable()) }()

	kvDB, err := database.NewKVDatabase(
		storage.NewSettings(&storage.SSettings{
			FPath: cDBPath,
		}),
	)
	if err != nil {
		return err
	}
	defer kvDB.Close()

	return adapters.ConsumeProcessor(
		pCtx,
		newAdaptedConsumer(p.fServiceAddr, p.fSettings, kvDB),
		p.fStdfLogger,
		hlt_client.NewClient(
			hlt_client.NewBuilder(),
			hlt_client.NewRequester(
				fmt.Sprintf("http://%s", p.fHltAddr),
				&http.Client{Timeout: time.Minute},
				p.fSettings,
			),
		),
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
