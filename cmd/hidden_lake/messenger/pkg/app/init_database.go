package app

import (
	"fmt"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/database"
	"github.com/number571/go-peer/pkg/storage"

	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/settings"
)

func (p *sApp) initDatabase() error {
	sett := storage.NewSettings(&storage.SSettings{
		FPath:     fmt.Sprintf("%s/%s", p.fPathTo, hlm_settings.CPathDB),
		FWorkSize: p.fConfig.GetSettings().GetWorkSizeBits(),
		FPassword: p.fConfig.GetStorageKey(),
	})
	db, err := database.NewKeyValueDB(sett)
	if err != nil {
		return fmt.Errorf("open KV database: %w", err)
	}
	p.fDatabase = db
	return nil
}
