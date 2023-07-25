package main

import (
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
}

func initApp() (types.ICommand, error) {
	hlsApp, err := hls_app.InitApp(".")
	if err != nil {
		return nil, err
	}

	hltApp, err := hlt_app.InitApp(".")
	if err != nil {
		return nil, err
	}

	return &sApp{
		fHLS: hlsApp,
		fHLT: hltApp,
	}, nil
}

func (p *sApp) Run() error {
	if err := p.fHLS.Run(); err != nil {
		return err
	}
	return p.fHLT.Run()
}

func (p *sApp) Stop() error {
	return errors.AppendError(p.fHLS.Stop(), p.fHLT.Stop())
}
