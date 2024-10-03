package config

import (
	logger "github.com/number571/go-peer/internal/logger/std"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	net_message "github.com/number571/go-peer/pkg/network/message"
)

type IWrapper interface {
	GetConfig() IConfig
	GetEditor() IEditor
}

type IEditor interface {
	UpdateConnections([]string) error
	UpdateFriends(map[string]asymmetric.IPubKey) error
}

type IConfigSettings interface {
	message.ISettings
	net_message.ISettings

	GetFetchTimeoutMS() uint64
	GetQueuePeriodMS() uint64
	GetRandQueuePeriodMS() uint64
	GetRandMessageSizeBytes() uint64
}

type IConfig interface {
	GetSettings() IConfigSettings
	GetLogging() logger.ILogging
	GetAddress() IAddress
	GetFriends() map[string]asymmetric.IPubKey
	GetConnections() []string
	GetService(string) (IService, bool)
}

type IService interface {
	GetHost() string
}

type IAddress interface {
	GetTCP() string
	GetHTTP() string
	GetPPROF() string
}
