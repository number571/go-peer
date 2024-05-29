package handler

import (
	"github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/internal/config"
	"github.com/number571/go-peer/internal/language"
)

type sTemplate struct {
	FLanguage language.ILanguage
}

func getTemplate(pCfg config.IConfig) *sTemplate {
	return &sTemplate{
		FLanguage: pCfg.GetSettings().GetLanguage(),
	}
}
