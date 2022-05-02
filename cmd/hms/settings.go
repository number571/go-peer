package main

import (
	"github.com/number571/go-peer/cmd/hms/config"
	"github.com/number571/go-peer/cmd/hms/database"
	"github.com/number571/go-peer/settings"
)

const (
	cErrorNone = iota + 1
	cErrorMethod
	cErrorDecode
	cErrorLoad
	cErrorPush
	cErrorMessage
	cErrorWorkSize
)

var (
	gSettings settings.ISettings
	gConfig   config.IConfig
	gDB       database.IKeyValueDB
)
