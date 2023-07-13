package config

import (
	"github.com/number571/go-peer/internal/logger"
	"github.com/number571/go-peer/internal/settings"
)

type IConfig interface {
	settings.IConfigSettings

	GetLogging() logger.ILogging
	GetNetwork() string
	GetStorage() bool
	GetAddress() IAddress
	GetConnections() []string
	GetConsumers() []string
}

type IAddress interface {
	GetTCP() string
	GetHTTP() string
}
