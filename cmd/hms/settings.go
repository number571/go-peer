package main

import (
	"github.com/number571/go-peer/cmd/hms/config"
	"github.com/number571/go-peer/cmd/hms/database"
)

var (
	gConfig config.IConfig
	gDB     database.IKeyValueDB
)
