package config

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
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

func (p *sEditor) UpdateNetworkKey(pNetworkKey string) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	filepath := p.fConfig.fFilepath
	icfg, err := LoadConfig(filepath)
	if err != nil {
		return fmt.Errorf("load config (update connections): %w", err)
	}

	cfg := icfg.(*SConfig)
	cfg.FSettings.FNetworkKey = pNetworkKey
	if err := os.WriteFile(filepath, encoding.Serialize(cfg, true), 0o644); err != nil {
		return fmt.Errorf("write config (update connections): %w", err)
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
		return fmt.Errorf("load config (update connections): %w", err)
	}

	cfg := icfg.(*SConfig)
	cfg.FConnections = deleteDuplicateStrings(pConns)
	if err := os.WriteFile(filepath, encoding.Serialize(cfg, true), 0o644); err != nil {
		return fmt.Errorf("write config (update connections): %w", err)
	}

	p.fConfig.fMutex.Lock()
	defer p.fConfig.fMutex.Unlock()

	p.fConfig.FConnections = cfg.FConnections
	return nil
}

func (p *sEditor) UpdateFriends(pFriends map[string]asymmetric.IPubKey) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	for name, pubKey := range pFriends {
		if pubKey.GetSize() == p.fConfig.GetSettings().GetKeySizeBits() {
			continue
		}
		return fmt.Errorf("not supported key size for '%s'", name)
	}

	filepath := p.fConfig.fFilepath
	icfg, err := LoadConfig(filepath)
	if err != nil {
		return fmt.Errorf("load config (update friends): %w", err)
	}

	if hasDuplicatePubKeys(pFriends) {
		return errors.New("has duplicates public keys")
	}

	cfg := icfg.(*SConfig)
	cfg.fFriends = pFriends
	cfg.FFriends = pubKeysToStrings(pFriends)
	if err := os.WriteFile(filepath, encoding.Serialize(cfg, true), 0o644); err != nil {
		return fmt.Errorf("write config (update friends): %w", err)
	}

	p.fConfig.fMutex.Lock()
	defer p.fConfig.fMutex.Unlock()

	p.fConfig.fFriends = cfg.fFriends
	p.fConfig.FFriends = cfg.FFriends
	return nil
}

func pubKeysToStrings(pPubKeys map[string]asymmetric.IPubKey) map[string]string {
	result := make(map[string]string, len(pPubKeys))
	for name, pubKey := range pPubKeys {
		result[name] = pubKey.ToString()
	}
	return result
}

func hasDuplicatePubKeys(pPubKeys map[string]asymmetric.IPubKey) bool {
	mapping := make(map[string]struct{})
	for _, pubKey := range pPubKeys {
		pubStr := pubKey.GetAddress().ToString()
		if _, ok := mapping[pubStr]; ok {
			return true
		}
		mapping[pubStr] = struct{}{}
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
