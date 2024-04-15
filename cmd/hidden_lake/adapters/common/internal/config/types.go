package config

import (
	logger "github.com/number571/go-peer/internal/logger/std"
	net_message "github.com/number571/go-peer/pkg/network/message"
)

type IConfig interface {
	GetLogging() logger.ILogging
	GetSettings() IConfigSettings

	GetAddress() string
	GetConnection() IConnection
}

type IConfigSettings interface {
	net_message.ISettings
	GetWaitTimeMS() uint64
}

type IConnection interface {
	GetHLTHost() string
	GetSrvHost() string
}
