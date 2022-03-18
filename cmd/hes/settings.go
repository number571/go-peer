package main

import (
	"github.com/number571/go-peer/cmd/hes/config"
	"github.com/number571/go-peer/cmd/hes/database"
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
