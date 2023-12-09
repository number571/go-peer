package app

import (
	"fmt"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/storage"
	"github.com/number571/go-peer/pkg/storage/database"
)

func (p *sApp) initDatabase() error {
	db, err := database.NewKVDatabase(
		storage.NewSettings(&storage.SSettings{
			FPath: fmt.Sprintf("%s/%s", p.fPathTo, pkg_settings.CPathDB),
		}),
	)
	if err != nil {
		return fmt.Errorf("new key/value database: %w", err)
	}
	p.fNode.GetWrapperDB().Set(db)
	return nil
}
