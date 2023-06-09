package app

import (
	"fmt"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/storage/database"
)

func (p *sApp) initDatabase() error {
	db, err := database.NewKeyValueDB(
		database.NewSettings(&database.SSettings{
			FPath:      fmt.Sprintf("%s/%s", p.fPathTo, pkg_settings.CPathDB),
			FHashing:   false,
			FCipherKey: []byte("_"),
		}),
	)
	if err != nil {
		return errors.WrapError(err, "new sqlitedb")
	}
	p.fNode.GetWrapperDB().Set(db)
	return nil
}
