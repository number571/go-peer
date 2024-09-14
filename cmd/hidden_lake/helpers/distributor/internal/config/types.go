package config

import (
	logger "github.com/number571/go-peer/internal/logger/std"
)

type IConfig interface {
	GetLogging() logger.ILogging
	GetSettings() IConfigSettings

	GetAddress() IAddress
	GetService(string) (string, bool)
}

type IConfigSettings interface{}

type IAddress interface {
	GetHTTP() string
	GetPPROF() string
}
