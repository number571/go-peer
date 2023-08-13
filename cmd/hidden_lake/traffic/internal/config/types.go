package config

import (
	logger "github.com/number571/go-peer/internal/logger/std"
	"github.com/number571/go-peer/internal/settings"
)

type IConfig interface {
	IConfigSettingsHLT

	GetLogging() logger.ILogging
	GetAddress() IAddress
	GetNetworkKey() string
	GetConnections() []string
	GetConsumers() []string
}

type IConfigSettingsHLT interface {
	IsValidHLT() bool
	settings.IConfigSettings

	GetMessagesCapacity() uint64
	GetQueuePeriodMS() uint64
	GetLimitVoidSizeBytes() uint64
}

type IAddress interface {
	GetTCP() string
	GetHTTP() string
}
