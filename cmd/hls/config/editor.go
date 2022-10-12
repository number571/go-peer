package config

import (
	"sync"

	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/encoding"
	"github.com/number571/go-peer/modules/filesystem"
)

var (
	_ IEditor = &sEditor{}
)

type sEditor struct {
	fMutex  sync.Mutex
	fConfig *IConfig
}

func NewEditor(cfg *IConfig) IEditor {
	if cfg == nil || *cfg == nil {
		return nil
	}
	return &sEditor{
		fConfig: cfg,
	}
}

func (edit *sEditor) UpdateConnections(conns []string) error {
	edit.fMutex.Lock()
	defer edit.fMutex.Unlock()

	filepath := (*edit.fConfig).(*sConfig).fFilepath
	bytes, err := filesystem.OpenFile(filepath).Read()
	if err != nil {
		return err
	}

	var cfg = new(sConfig)
	err = encoding.Deserialize(bytes, cfg)
	if err != nil {
		return err
	}

	cfg.FConnections = deleteDuplicateStrings(conns)
	err = filesystem.OpenFile(filepath).Write(encoding.Serialize(cfg))
	if err != nil {
		return err
	}

	(*edit.fConfig).(*sConfig).FConnections = cfg.FConnections
	return nil
}

func (edit *sEditor) UpdateFriends(friends []asymmetric.IPubKey) error {
	edit.fMutex.Lock()
	defer edit.fMutex.Unlock()

	filepath := (*edit.fConfig).(*sConfig).fFilepath
	bytes, err := filesystem.OpenFile(filepath).Read()
	if err != nil {
		return err
	}

	var cfg = new(sConfig)
	err = encoding.Deserialize(bytes, cfg)
	if err != nil {
		return err
	}

	cfg.fFriends = deleteDuplicatePubKeys(friends)
	cfg.FFriends = pubKeysToStrings(friends)
	err = filesystem.OpenFile(filepath).Write(encoding.Serialize(cfg))
	if err != nil {
		return err
	}

	(*edit.fConfig).(*sConfig).fFriends = cfg.fFriends
	(*edit.fConfig).(*sConfig).FFriends = cfg.FFriends
	return nil
}

func pubKeysToStrings(pubKeys []asymmetric.IPubKey) []string {
	result := make([]string, 0, len(pubKeys))
	for _, pk := range pubKeys {
		result = append(result, pk.String())
	}
	return result
}

func deleteDuplicatePubKeys(pubKeys []asymmetric.IPubKey) []asymmetric.IPubKey {
	result := make([]asymmetric.IPubKey, 0, len(pubKeys))
	mapping := make(map[string]struct{})
	for _, pk := range pubKeys {
		if _, ok := mapping[pk.Address().String()]; ok {
			continue
		}
		mapping[pk.Address().String()] = struct{}{}
		result = append(result, pk)
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
