package app

import (
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/config"
	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/settings"
	"github.com/number571/go-peer/pkg/storage"
)

func initCryptoStorage(cfg config.IConfig) (storage.IKeyValueStorage, error) {
	return storage.NewCryptoStorage(
		hlm_settings.CPathSTG,
		[]byte(cfg.GetStorageKey()),
		hlm_settings.CWorkForKeys,
	)
}
