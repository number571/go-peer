package main

import (
	hls_app "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/app"
	hlt_app "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/app"
	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/go-peer/pkg/utils"
)

var (
	_ types.IApp = &sApp{}
)

type sApp struct {
	fHLS types.IApp
	fHLT types.IApp
}

func initApp(pPath, pKey string) (types.IApp, error) {
	hlsApp, err := hls_app.InitApp(pPath, pKey)
	if err != nil {
		return nil, err
	}

	hltApp, err := hlt_app.InitApp(pPath)
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
	return utils.MergeErrors(p.fHLS.Stop(), p.fHLT.Stop())
}
