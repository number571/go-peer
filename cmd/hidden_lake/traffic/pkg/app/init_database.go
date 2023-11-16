package app

import (
	"fmt"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	"github.com/number571/go-peer/pkg/errors"

	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
)

func (p *sApp) initDatabase() error {
	sett := database.NewSettings(&database.SSettings{
		FPath:             fmt.Sprintf("%s/%s", p.fPathTo, hlt_settings.CPathDB),
		FNetworkKey:       p.fConfig.GetNetworkKey(),
		FMessagesCapacity: p.fConfig.GetSettings().GetMessagesCapacity(),
		FWorkSizeBits:     p.fConfig.GetSettings().GetWorkSizeBits(),
	})

	if !p.fConfig.GetIsStorage() {
		p.fWrapperDB.Set(database.NewVoidDatabase(sett))
		return nil
	}

	db, err := database.NewDatabase(sett)
	if err != nil {
		return errors.WrapError(err, "init database")
	}

	p.fWrapperDB.Set(db)
	return nil
}
