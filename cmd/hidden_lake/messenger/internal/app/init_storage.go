package app

import (
	"fmt"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/config"
	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/settings"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/storage"
)

func initCryptoStorage(pCfg config.IConfig, pPathTo string) (storage.IKeyValueStorage, error) {
	stg, err := storage.NewCryptoStorage(storage.NewSettings(&storage.SSettings{
		FPath:      fmt.Sprintf("%s/%s", pPathTo, hlm_settings.CPathSTG),
		FWorkSize:  hlm_settings.CWorkForKeys,
		FCipherKey: []byte(pCfg.GetStorageKey()),
	}))
	if err != nil {
		return nil, errors.WrapError(err, "new crypto storage")
	}
	return stg, nil
}
