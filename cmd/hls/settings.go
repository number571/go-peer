package main

import (
	"net/http"

	"github.com/number571/go-peer/cmd/hls/config"
	"github.com/number571/go-peer/logger"
	"github.com/number571/go-peer/network/anonymity"
)

var (
	gLogger     logger.ILogger
	gNode       anonymity.INode
	gConfig     config.IConfig
	gServerHTTP *http.Server
)
