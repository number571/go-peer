package app

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/number571/go-peer/cmd/hidden_lake/loader/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/loader/pkg/settings"
	std_logger "github.com/number571/go-peer/internal/logger/std"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/types"
)

var (
	_ types.ICommand = &sApp{}
)

type sApp struct {
	fIsRun      bool
	fMutex      sync.Mutex
	fStdfLogger logger.ILogger
	fCancel     context.CancelFunc
	fConfig     config.IConfig
}

func NewApp(
	pCfg config.IConfig,
	pPathTo string,
) types.ICommand {
	stdfLogger := std_logger.NewStdLogger(
		pCfg.GetLogging(),
		std_logger.GetLogFunc(),
	)

	return &sApp{
		fConfig:     pCfg,
		fStdfLogger: stdfLogger,
	}
}

func (p *sApp) Run() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if p.fIsRun {
		return errors.New("application already running")
	}
	p.fIsRun = true

	ctx := context.Background()
	ctxWithCancel, cancelFunction := context.WithCancel(ctx)

	p.fCancel = cancelFunction
	p.transferMessages(ctxWithCancel)

	p.fStdfLogger.PushInfo(fmt.Sprintf("%s is running...", settings.CServiceName))
	return nil
}

func (p *sApp) Stop() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if !p.fIsRun {
		return errors.New("application already stopped or not started")
	}
	p.fIsRun = false
	p.fCancel()

	p.fStdfLogger.PushInfo(fmt.Sprintf("%s is shutting down...", settings.CServiceName))
	return nil
}
