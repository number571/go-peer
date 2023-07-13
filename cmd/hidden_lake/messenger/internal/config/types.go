package config

import (
	"github.com/number571/go-peer/internal/logger"
	"github.com/number571/go-peer/internal/settings"
)

type IConfig interface {
	settings.IConfigSettings

	GetAddress() IAddress
	GetLogging() logger.ILogging
	GetConnection() IConnection
	GetStorageKey() string
}

type IAddress interface {
	GetInterface() string
	GetIncoming() string
}

type IConnection interface {
	GetService() string
	GetTraffic() string
}
