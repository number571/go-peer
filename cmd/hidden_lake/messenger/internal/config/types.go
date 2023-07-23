package config

import (
	"github.com/number571/go-peer/internal/logger"
	"github.com/number571/go-peer/internal/settings"
)

type IConfig interface {
	IConfigSettingsHLM

	GetAddress() IAddress
	GetLogging() logger.ILogging
	GetConnection() IConnection
	GetStorageKey() string
}

type IConfigSettingsHLM interface {
	IsValidHLM() bool
	settings.IConfigSettings

	GetKeySizeBits() uint64
	GetMessagesCapacity() uint64
}

type IAddress interface {
	GetInterface() string
	GetIncoming() string
}

type IConnection interface {
	GetService() string
	GetTraffic() string
}
