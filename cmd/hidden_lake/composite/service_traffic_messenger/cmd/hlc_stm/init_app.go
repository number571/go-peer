package main

import (
	"context"
	"errors"
	"sync"

	hlm_app "github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/app"
	hls_app "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/app"
	hlt_app "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/app"
	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/go-peer/pkg/utils"
)

var (
	_ types.IRunner = &sApp{}
)

type sApp struct {
	fMutex sync.Mutex
	fIsRun bool
	fHLS   types.IRunner
	fHLT   types.IRunner
	fHLM   types.IRunner
}

func initApp(pPath, pKey string) (types.IRunner, error) {
	hlsApp, err := hls_app.InitApp(pPath, pKey)
	if err != nil {
		return nil, err
	}

	hltApp, err := hlt_app.InitApp(pPath)
	if err != nil {
		return nil, err
	}

	hlmApp, err := hlm_app.InitApp(pPath)
	if err != nil {
		return nil, err
	}

	return &sApp{
		fHLS: hlsApp,
		fHLT: hltApp,
		fHLM: hlmApp,
	}, nil
}

func (p *sApp) Run(pCtx context.Context) error {
	err := func() error {
		p.fMutex.Lock()
		defer p.fMutex.Unlock()

		if p.fIsRun {
			return errors.New("application already running")
		}

		p.fIsRun = true
		return nil
	}()
	if err != nil {
		return err
	}
	defer func() {
		p.fMutex.Lock()
		p.fIsRun = false
		p.fMutex.Unlock()
	}()

	var (
		hlsErr = make(chan error)
		hltErr = make(chan error)
		hlmErr = make(chan error)
	)

	go func() {
		hlsErr <- p.fHLS.Run(pCtx)
	}()
	go func() {
		hltErr <- p.fHLT.Run(pCtx)
	}()
	go func() {
		hlmErr <- p.fHLM.Run(pCtx)
	}()

	return utils.MergeErrors(<-hlsErr, <-hltErr, <-hlmErr)
}
