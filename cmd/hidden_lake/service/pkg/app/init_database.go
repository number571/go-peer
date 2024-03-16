package app

import (
	"fmt"
	"path/filepath"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/database"
)

func (p *sApp) initDatabase() error {
	db, err := database.NewKVDatabase(
		database.NewSettings(&database.SSettings{
			FPath: filepath.Join(p.fPathTo, pkg_settings.CPathDB),
		}),
	)
	if err != nil {
		return fmt.Errorf("new key/value database: %w", err)
	}
	p.fNode.GetDBWrapper().Set(db)
	return nil
}
