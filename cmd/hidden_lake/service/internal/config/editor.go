package config

import (
	"os"
	"sync"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/utils"
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

func (p *sEditor) UpdateNetworkKey(pNetworkKey string) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	filepath := p.fConfig.fFilepath
	icfg, err := LoadConfig(filepath)
	if err != nil {
		return utils.MergeErrors(ErrLoadConfig, err)
	}

	cfg := icfg.(*SConfig)
	cfg.FSettings.FNetworkKey = pNetworkKey
	if err := os.WriteFile(filepath, encoding.SerializeYAML(cfg), 0o600); err != nil {
		return utils.MergeErrors(ErrWriteConfig, err)
	}

	p.fConfig.FSettings.fMutex.Lock()
	defer p.fConfig.FSettings.fMutex.Unlock()

	p.fConfig.FSettings.FNetworkKey = cfg.FSettings.FNetworkKey
	return nil
}

func (p *sEditor) UpdateConnections(pConns []string) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	filepath := p.fConfig.fFilepath
	icfg, err := LoadConfig(filepath)
	if err != nil {
		return utils.MergeErrors(ErrLoadConfig, err)
	}

	cfg := icfg.(*SConfig)
	cfg.FConnections = deleteDuplicateStrings(pConns)
	if err := os.WriteFile(filepath, encoding.SerializeYAML(cfg), 0o600); err != nil {
		return utils.MergeErrors(ErrWriteConfig, err)
	}

	p.fConfig.fMutex.Lock()
	defer p.fConfig.fMutex.Unlock()

	p.fConfig.FConnections = cfg.FConnections
	return nil
}

func (p *sEditor) UpdateFriends(pFriends map[string]string) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	filepath := p.fConfig.fFilepath
	icfg, err := LoadConfig(filepath)
	if err != nil {
		return utils.MergeErrors(ErrLoadConfig, err)
	}

	if hasDuplicateKeys(pFriends) {
		return ErrDuplicatePublicKey
	}

	cfg := icfg.(*SConfig)
	cfg.FFriends = pFriends
	if err := os.WriteFile(filepath, encoding.SerializeYAML(cfg), 0o600); err != nil {
		return utils.MergeErrors(ErrWriteConfig, err)
	}

	p.fConfig.fMutex.Lock()
	defer p.fConfig.fMutex.Unlock()

	p.fConfig.FFriends = cfg.FFriends
	return nil
}

func hasDuplicateKeys(pKeys map[string]string) bool {
	mapping := make(map[string]struct{})
	for _, key := range pKeys {
		if _, ok := mapping[key]; ok {
			return true
		}
		mapping[key] = struct{}{}
	}
	return false
}

func deleteDuplicateStrings(pStrs []string) []string {
	result := make([]string, 0, len(pStrs))
	mapping := make(map[string]struct{})
	for _, s := range pStrs {
		if _, ok := mapping[s]; ok {
			continue
		}
		mapping[s] = struct{}{}
		result = append(result, s)
	}
	return result
}
