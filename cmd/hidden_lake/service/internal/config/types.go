package config

import (
	logger "github.com/number571/go-peer/internal/logger/std"
	"github.com/number571/go-peer/internal/settings"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

type IWrapper interface {
	GetConfig() IConfig
	GetEditor() IEditor
}

type IEditor interface {
	UpdateNetworkKey(string) error
	UpdateConnections([]string) error
	UpdateFriends(map[string]asymmetric.IPubKey) error
}

type IConfig interface {
	IConfigSettingsHLS

	GetLogging() logger.ILogging
	GetAddress() IAddress
	GetNetworkKey() string
	GetConnections() []string
	GetFriends() map[string]asymmetric.IPubKey
	GetService(string) (string, bool)
}

type IConfigSettingsHLS interface {
	IsValidHLS() bool
	settings.IConfigSettings

	GetKeySizeBits() uint64
	GetQueuePeriodMS() uint64
	GetLimitVoidSizeBytes() uint64
}

type IAddress interface {
	GetTCP() string
	GetHTTP() string
}
