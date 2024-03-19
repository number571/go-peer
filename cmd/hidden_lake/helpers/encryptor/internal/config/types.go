package config

import (
	logger "github.com/number571/go-peer/internal/logger/std"
	"github.com/number571/go-peer/pkg/client/message"
	net_message "github.com/number571/go-peer/pkg/network/message"
)

type IConfig interface {
	GetLogging() logger.ILogging
	GetSettings() IConfigSettings

	GetAddress() IAddress
}

type IConfigSettings interface {
	message.ISettings
	net_message.ISettings

	GetLimitVoidSizeBytes() uint64
}

type IAddress interface {
	GetHTTP() string
	GetPPROF() string
}
