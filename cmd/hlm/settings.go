package main

import (
	"net/http"

	"github.com/number571/go-peer/cmd/hlm/config"
	"github.com/number571/go-peer/cmd/hls/hlc"
	"github.com/number571/go-peer/modules/action"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/logger"
)

var (
	gLogger        logger.ILogger
	gActions       action.IActions
	gConfig        config.IConfig
	gClient        hlc.IClient
	gServerHTTP    *http.Server
	gChannelPubKey asymmetric.IPubKey
)
