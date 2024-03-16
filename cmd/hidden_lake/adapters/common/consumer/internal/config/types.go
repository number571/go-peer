package config

import (
	logger "github.com/number571/go-peer/internal/logger/std"
	net_message "github.com/number571/go-peer/pkg/network/message"
)

type IConfig interface {
	GetLogging() logger.ILogging
	GetSettings() IConfigSettings

	GetConnection() IConnection
}

type IConfigSettings interface {
	net_message.ISettings
}

type IConnection interface {
	GetHLTHost() string
	GetSrvHost() string
}
