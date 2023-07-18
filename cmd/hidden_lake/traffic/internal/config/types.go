package config

import (
	"github.com/number571/go-peer/internal/logger"
	"github.com/number571/go-peer/internal/settings"
)

type IConfig interface {
	IConfigSettingsHLT

	GetLogging() logger.ILogging
	GetNetwork() string
	GetAddress() IAddress
	GetConnections() []string
	GetConsumers() []string
}

type IConfigSettingsHLT interface {
	IsValidHLT() bool
	settings.IConfigSettings

	GetMessagesCapacity() uint64
}

type IAddress interface {
	GetTCP() string
	GetHTTP() string
}
