package config

import "github.com/number571/go-peer/pkg/logger"

type IConfig interface {
	GetLogging() logger.ILogging
	GetNetwork() string
	GetAddress() string
	GetConnection() string
	GetConsumers() []string
}
