package config

import (
	logger "github.com/number571/go-peer/internal/logger/std"
	"github.com/number571/go-peer/pkg/client/message"
)

type IConfig interface {
	GetSettings() IConfigSettings
	GetLogging() logger.ILogging
	GetAddress() IAddress
	GetNetworkKey() string
	GetConnections() []string
	GetConsumers() []string
}

type IConfigSettings interface {
	message.ISettings

	GetQueuePeriodMS() uint64
	GetMessagesCapacity() uint64
	GetLimitVoidSizeBytes() uint64
}

type IAddress interface {
	GetTCP() string
	GetHTTP() string
	GetPPROF() string
}
