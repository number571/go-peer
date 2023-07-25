package main

import (
	hlm_app "github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/app"
	hls_app "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/app"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/types"
)

var (
	_ types.ICommand = &sApp{}
)

type sApp struct {
	fHLS types.ICommand
	fHLM types.ICommand
}

func initApp() (types.ICommand, error) {
	hlsApp, err := hls_app.InitApp(".")
	if err != nil {
		return nil, err
	}

	hlmApp, err := hlm_app.InitApp(".")
	if err != nil {
		return nil, err
	}

	return &sApp{
		fHLS: hlsApp,
		fHLM: hlmApp,
	}, nil
}

func (p *sApp) Run() error {
	if err := p.fHLS.Run(); err != nil {
		return err
	}
	return p.fHLM.Run()
}

func (p *sApp) Stop() error {
	err := p.fHLM.Stop()
	return errors.AppendError(p.fHLS.Stop(), err)
}
