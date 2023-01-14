package config

import (
	"fmt"
	"sync"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/filesystem"
)

var (
	_ IEditor = &sEditor{}
)

type sEditor struct {
	fMutex  sync.Mutex
	fConfig *SConfig
}

func newEditor(cfg IConfig) IEditor {
	if cfg == nil {
		return nil
	}
	v, ok := cfg.(*SConfig)
	if !ok {
		return nil
	}
	return &sEditor{
		fConfig: v,
	}
}

func (edit *sEditor) UpdateConnections(conns []string) error {
	edit.fMutex.Lock()
	defer edit.fMutex.Unlock()

	filepath := edit.fConfig.fFilepath
	icfg, err := LoadConfig(filepath)
	if err != nil {
		return err
	}

	cfg := icfg.(*SConfig)
	cfg.FConnections = deleteDuplicateStrings(conns)
	err = filesystem.OpenFile(filepath).Write(encoding.Serialize(cfg))
	if err != nil {
		return err
	}

	edit.fConfig.fMutex.Lock()
	defer edit.fConfig.fMutex.Unlock()

	edit.fConfig.FConnections = cfg.FConnections
	return nil
}

func (edit *sEditor) UpdateFriends(friends map[string]asymmetric.IPubKey) error {
	edit.fMutex.Lock()
	defer edit.fMutex.Unlock()

	for name, pubKey := range friends {
		if pubKey.Size() == pkg_settings.CAKeySize {
			continue
		}
		return fmt.Errorf("not supported key size for '%s'", name)
	}

	filepath := edit.fConfig.fFilepath
	icfg, err := LoadConfig(filepath)
	if err != nil {
		return err
	}

	cfg := icfg.(*SConfig)
	cfg.fFriends = deleteDuplicatePubKeys(friends)
	cfg.FFriends = pubKeysToStrings(friends)
	err = filesystem.OpenFile(filepath).Write(encoding.Serialize(cfg))
	if err != nil {
		return err
	}

	edit.fConfig.fMutex.Lock()
	defer edit.fConfig.fMutex.Unlock()

	edit.fConfig.fFriends = cfg.fFriends
	edit.fConfig.FFriends = cfg.FFriends
	return nil
}

func pubKeysToStrings(pubKeys map[string]asymmetric.IPubKey) map[string]string {
	result := make(map[string]string, len(pubKeys))
	for name, pubKey := range pubKeys {
		result[name] = pubKey.String()
	}
	return result
}

func deleteDuplicatePubKeys(pubKeys map[string]asymmetric.IPubKey) map[string]asymmetric.IPubKey {
	result := make(map[string]asymmetric.IPubKey, len(pubKeys))
	mapping := make(map[string]struct{})
	for name, pubKey := range pubKeys {
		pubStr := pubKey.Address().String()
		if _, ok := mapping[pubStr]; ok {
			continue
		}
		mapping[pubStr] = struct{}{}
		result[name] = pubKey
	}
	return result
}

func deleteDuplicateStrings(strs []string) []string {
	result := make([]string, 0, len(strs))
	mapping := make(map[string]struct{})
	for _, s := range strs {
		if _, ok := mapping[s]; ok {
			continue
		}
		mapping[s] = struct{}{}
		result = append(result, s)
	}
	return result
}
