package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/utils"
	"github.com/number571/go-peer/pkg/encoding"
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
		panic("cfg = nil")
	}
	v, ok := pCfg.(*SConfig)
	if !ok {
		panic("cfg is invalid")
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
		return fmt.Errorf("load config (update language): %w", err)
	}

	cfg := icfg.(*SConfig)
	cfg.FLanguage = utils.FromILanguage(pLang)
	if err := os.WriteFile(filepath, encoding.SerializeYAML(cfg), 0o644); err != nil {
		return fmt.Errorf("write config (update language): %w", err)
	}

	p.fConfig.fMutex.Lock()
	defer p.fConfig.fMutex.Unlock()

	p.fConfig.fLanguage = pLang
	return nil
}

func (p *sEditor) UpdateSecretKeys(pSecretKeys map[string]string) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	filepath := p.fConfig.fFilepath
	icfg, err := LoadConfig(filepath)
	if err != nil {
		return fmt.Errorf("load config (update friends): %w", err)
	}

	cfg := icfg.(*SConfig)
	cfg.FSecretKeys = pSecretKeys
	if err := os.WriteFile(filepath, encoding.SerializeYAML(cfg), 0o644); err != nil {
		return fmt.Errorf("write config (update friends): %w", err)
	}

	p.fConfig.fMutex.Lock()
	defer p.fConfig.fMutex.Unlock()

	p.fConfig.FSecretKeys = cfg.FSecretKeys
	return nil
}
