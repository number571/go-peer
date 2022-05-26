package main

import (
	"github.com/number571/go-peer/cmd/hms/config"
	"github.com/number571/go-peer/cmd/hms/database"
	"github.com/number571/go-peer/settings"
)

var (
	gSettings settings.ISettings
	gConfig   config.IConfig
	gDB       database.IKeyValueDB
)
