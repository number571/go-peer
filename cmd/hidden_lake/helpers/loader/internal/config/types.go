package config

import (
	logger "github.com/number571/go-peer/internal/logger/std"
	net_message "github.com/number571/go-peer/pkg/network/message"
)

type IConfig interface {
	GetLogging() logger.ILogging
	GetSettings() IConfigSettings

	GetAddress() IAddress
	GetProducers() []string
	GetConsumers() []string
}

type IConfigSettings interface {
	net_message.ISettings

	GetTimestampWindowS() uint64
	GetMessagesCapacity() uint64
}

type IAddress interface {
	GetHTTP() string
	GetPPROF() string
}
