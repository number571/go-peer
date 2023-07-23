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
		FCapacity:         p.fConfig.GetMessagesCapacity(),
		FMessageSizeBytes: p.fConfig.GetMessageSizeBytes(),
		FWorkSizeBits:     p.fConfig.GetWorkSizeBits(),
	})

	if p.fConfig.GetMessagesCapacity() == 0 {
		p.fWrapperDB.Set(database.NewVoidKeyValueDB(sett))
		return nil
	}

	db, err := database.NewKeyValueDB(sett)
	if err != nil {
		return errors.WrapError(err, "init database")
	}

	p.fWrapperDB.Set(db)
	return nil
}
