package main

import (
	"github.com/number571/go-peer/cmd/hls/config"
	"github.com/number571/go-peer/cmd/hls/database"
	"github.com/number571/go-peer/cmd/hls/logger"
	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/network"
)

var (
	gPPrivKey asymmetric.IPrivKey
	gLogger   logger.ILogger
	gNode     network.INode
	gConfig   config.IConfig
	gDB       database.IKeyValueDB
)
