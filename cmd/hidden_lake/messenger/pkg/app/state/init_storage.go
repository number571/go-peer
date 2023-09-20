package state

import (
	"fmt"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/config"
	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/settings"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/storage"
)

func initCryptoStorage(pCfg config.IConfig, pPathTo string) (storage.IKVStorage, error) {
	storageKey := pCfg.GetStorageKey()
	if storageKey == "" {
		return nil, errors.NewError("storage key is nil")
	}
	sett := storage.NewSettings(&storage.SSettings{
		FPath:     fmt.Sprintf("%s/%s", pPathTo, hlm_settings.CPathSTG),
		FWorkSize: hlm_settings.CWorkForKeys,
		FPassword: storageKey,
	})
	stg, err := storage.NewCryptoStorage(sett)
	if err != nil {
		return nil, errors.WrapError(err, "new crypto storage")
	}
	return stg, nil
}
