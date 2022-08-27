package main

import (
	"net/http"

	"github.com/number571/go-peer/cmd/hls/config"
	"github.com/number571/go-peer/modules/logger"
	"github.com/number571/go-peer/modules/network/anonymity"
	"github.com/number571/go-peer/modules/network/conn_keeper"
)

var (
	gLogger     logger.ILogger
	gNode       anonymity.INode
	gConfig     config.IConfig
	gConnKeeper conn_keeper.IConnKeeper
	gServerHTTP *http.Server
)
