package config

import (
	logger "github.com/number571/go-peer/internal/logger/std"
)

type IConfig interface {
	GetSettings() IConfigSettings
	GetAddress() IAddress
	GetLogging() logger.ILogging
}

type IConfigSettings interface {
	GetExecTimeoutMS() uint64
}

type IAddress interface {
	GetIncoming() string
	GetPPROF() string
}
