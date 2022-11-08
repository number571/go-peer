package main

import (
	"net/http"

	"github.com/number571/go-peer/cmd/hls/config"
	"github.com/number571/go-peer/modules/logger"
	"github.com/number571/go-peer/modules/network"
	"github.com/number571/go-peer/modules/network/anonymity"
	"github.com/number571/go-peer/modules/network/conn_keeper"
	"github.com/number571/go-peer/modules/storage/database"
)

var (
	gLogger      logger.ILogger
	gNode        anonymity.INode
	gConfig      config.IConfig
	gEditor      config.IEditor
	gConnKeeper  conn_keeper.IConnKeeper
	gServerHTTP  *http.Server
	gLevelDB     database.IKeyValueDB
	gNetworkNode network.INode
)
