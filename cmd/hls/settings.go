package main

import (
	"github.com/number571/go-peer/cmd/hls/config"
	"github.com/number571/go-peer/cmd/hls/database"
	"github.com/number571/go-peer/cmd/hls/logger"
	"github.com/number571/go-peer/network"
)

const (
	cErrorNone = iota + 1
	cErrorMethod
	cErrorDecodeRequest
	cErrorDecodePubKey
	cErrorResponseMessage
)

const (
	cPatternHLS = "hidden-lake-service"
)

const (
	cAKeySize = 4096
	cProto    = "http"
)

var (
	gLogger logger.ILogger
	gNode   network.INode
	gConfig config.IConfig
	gDB     database.IKeyValueDB
)
