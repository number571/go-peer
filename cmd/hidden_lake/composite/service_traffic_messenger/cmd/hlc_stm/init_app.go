package main

import (
	hlm_app "github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/app"
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
	fHLM types.IApp
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
	return utils.MergeErrors(p.fHLS.Stop(), p.fHLT.Stop(), p.fHLS.Stop())
}
