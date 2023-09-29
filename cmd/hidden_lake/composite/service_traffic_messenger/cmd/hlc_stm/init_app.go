package main

import (
	hlm_app "github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/app"
	hls_app "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/app"
	hlt_app "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/app"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/types"
)

var (
	_ types.ICommand = &sApp{}
)

type sApp struct {
	fHLS types.ICommand
	fHLT types.ICommand
	fHLM types.ICommand
}

func initApp(pPath string) (types.ICommand, error) {
	hlsApp, err := hls_app.InitApp(pPath)
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

func (p *sApp) Run() error {
	if err := p.fHLS.Run(); err != nil {
		return err
	}
	if err := p.fHLT.Run(); err != nil {
		return err
	}
	return p.fHLM.Run()
}

func (p *sApp) Stop() error {
	errT := p.fHLM.Stop()
	err := errors.AppendError(p.fHLS.Stop(), errT)
	return errors.AppendError(p.fHLT.Stop(), err)
}
