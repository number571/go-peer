package main

import (
	"github.com/number571/go-peer/cmd/hls/config"
	"github.com/number571/go-peer/local"
)

const (
	cAKeySize   = 4096
	cPatternHLS = "hidden-lake-service"
)

var (
	gClient local.IClient
	gConfig config.IConfig
)
