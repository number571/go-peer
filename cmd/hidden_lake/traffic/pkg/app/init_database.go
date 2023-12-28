package app

import (
	"fmt"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"

	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
)

func (p *sApp) initDatabase() error {
	cfgSettings := p.fConfig.GetSettings()
	sett := database.NewSettings(&database.SSettings{
		FPath:             fmt.Sprintf("%s/%s", p.fPathTo, hlt_settings.CPathDB),
		FNetworkKey:       cfgSettings.GetNetworkKey(),
		FWorkSizeBits:     cfgSettings.GetWorkSizeBits(),
		FMessagesCapacity: cfgSettings.GetMessagesCapacity(),
	})

	var (
		db  database.IDatabase
		err error
	)

	switch {
	case p.fConfig.GetStorage():
		db, err = database.NewDatabase(sett)
	default:
		db, err = database.NewInMemoryDatabase(sett)
	}
	if err != nil {
		return fmt.Errorf("init database: %w", err)
	}

	p.fDBWrapper.Set(db)
	return nil
}
