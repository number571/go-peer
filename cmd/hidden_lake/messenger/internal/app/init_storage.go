package app

import (
	"fmt"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/config"
	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/settings"
	"github.com/number571/go-peer/pkg/storage"
)

func initCryptoStorage(pCfg config.IConfig, pPathTo string) (storage.IKeyValueStorage, error) {
	return storage.NewCryptoStorage(
		fmt.Sprintf("%s/%s", pPathTo, hlm_settings.CPathSTG),
		[]byte(pCfg.GetStorageKey()),
		hlm_settings.CWorkForKeys,
	)
}
