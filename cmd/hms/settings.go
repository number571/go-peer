package main

import (
	"github.com/number571/go-peer/cmd/hms/config"
	"github.com/number571/go-peer/cmd/hms/database"
)

const (
	cErrorNone = iota + 1
	cErrorMethod
	cErrorDecode
	cErrorSize
	cErrorLoad
	cErrorPush
	cErrorMessage
)

var (
	gDB     database.IKeyValueDB
	gConfig config.IConfig
)
