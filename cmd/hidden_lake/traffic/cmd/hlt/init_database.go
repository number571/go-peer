package main

import (
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
)

func initDatabase() database.IKeyValueDB {
	return database.NewKeyValueDB(
		database.NewSettings(&database.SSettings{
			FPath:          hlt_settings.CPathDB,
			FLimitMessages: hlt_settings.CLimitMessages,
			FMessageSize:   hlt_settings.CMessageSize,
			FWorkSize:      hlt_settings.CWorkSize,
		}),
	)
}
