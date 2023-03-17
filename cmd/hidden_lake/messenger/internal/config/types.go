package config

import "github.com/number571/go-peer/pkg/logger"

type IConfig interface {
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
