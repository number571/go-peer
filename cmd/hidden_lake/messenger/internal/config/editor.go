package config

import (
	"sync"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/utils"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/file_system"
)

var (
	_ IEditor = &sEditor{}
)

type sEditor struct {
	fMutex  sync.Mutex
	fConfig *SConfig
}

func newEditor(pCfg IConfig) IEditor {
	if pCfg == nil {
		return nil
	}
	v, ok := pCfg.(*SConfig)
	if !ok {
		return nil
	}
	return &sEditor{
		fConfig: v,
	}
}

func (p *sEditor) UpdateLanguage(pLang utils.ILanguage) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	filepath := p.fConfig.fFilepath
	icfg, err := LoadConfig(filepath)
	if err != nil {
		return errors.WrapError(err, "load config (update language)")
	}

	cfg := icfg.(*SConfig)
	cfg.FLanguage = utils.FromILanguage(pLang)
	err = file_system.OpenFile(filepath).Write(encoding.Serialize(cfg, true))
	if err != nil {
		return errors.WrapError(err, "write config (update language)")
	}

	p.fConfig.fMutex.Lock()
	defer p.fConfig.fMutex.Unlock()

	p.fConfig.fLanguage = pLang
	return nil
}
