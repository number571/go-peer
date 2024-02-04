package app

import (
	"fmt"
	"path/filepath"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/database"
	"github.com/number571/go-peer/pkg/storage"

	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/pkg/settings"
)

func (p *sApp) initDatabase() error {
	sett := storage.NewSettings(&storage.SSettings{
		FPath:     filepath.Join(p.fPathTo, hlm_settings.CPathDB),
		FWorkSize: p.fConfig.GetSettings().GetWorkSizeBits(),
		FPassword: p.fConfig.GetSettings().GetStorageKey(),
	})
	db, err := database.NewKeyValueDB(sett)
	if err != nil {
		return fmt.Errorf("open KV database: %w", err)
	}
	p.fDatabase = db
	return nil
}
