package app

import (
	"fmt"
	"path/filepath"

	hlm_database "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/database"
	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/pkg/settings"
)

func (p *sApp) initDatabase() error {
	path := filepath.Join(p.fPathTo, hlm_settings.CPathDB)
	db, err := hlm_database.NewKeyValueDB(path)
	if err != nil {
		return fmt.Errorf("open KV database: %w", err)
	}
	p.fDatabase = db
	return nil
}
