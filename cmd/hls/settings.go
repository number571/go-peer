package main

import (
	"github.com/number571/go-peer/cmd/hls/config"
	"github.com/number571/go-peer/cmd/hls/logger"
	"github.com/number571/go-peer/netanon"
)

var (
	gLogger logger.ILogger
	gNode   netanon.INode
	gConfig config.IConfig
)
