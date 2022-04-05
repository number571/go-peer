package main

import (
	"github.com/number571/go-peer/cmd/hls/config"
	"github.com/number571/go-peer/cmd/hls/database"
	"github.com/number571/go-peer/cmd/hls/logger"
	"github.com/number571/go-peer/local"
)

const (
	cAKeySize   = 4096
	cPatternHLS = "hidden-lake-service"
	cProto      = "http"
)

var (
	gLogger logger.ILogger
	gClient local.IClient
	gConfig config.IConfig
	gDB     database.IKeyValueDB
)
