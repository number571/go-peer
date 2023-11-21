package config

import logger "github.com/number571/go-peer/internal/logger/std"

type IConfig interface {
	GetLogging() logger.ILogging
	GetSettings() IConfigSettings

	GetAddress() IAddress
	GetProducers() []string
	GetConsumers() []string
}

type IConfigSettings interface {
	GetNetworkKey() string
	GetWorkSizeBits() uint64
	GetMessagesCapacity() uint64
}

type IAddress interface {
	GetHTTP() string
	GetPPROF() string
}
