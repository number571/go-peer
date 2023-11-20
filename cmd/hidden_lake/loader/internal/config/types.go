package config

import logger "github.com/number571/go-peer/internal/logger/std"

type IConfig interface {
	GetLogging() logger.ILogging
	GetSettings() IConfigSettings

	GetProducers() []string
	GetConsumers() []string
}

type IConfigSettings interface {
	GetNetworkKey() string
	GetWorkSizeBits() uint64
	GetMessagesCapacity() uint64
}
