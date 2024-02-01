package config

import (
	logger "github.com/number571/go-peer/internal/logger/std"
)

type IConfig interface {
	GetLogging() logger.ILogging
	GetServices() []string
}
