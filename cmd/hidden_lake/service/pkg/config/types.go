package config

import (
	"github.com/number571/go-peer/internal/logger"
	"github.com/number571/go-peer/internal/settings"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

type IWrapper interface {
	GetConfig() IConfig
	GetEditor() IEditor
}

type IEditor interface {
	UpdateConnections([]string) error
	UpdateFriends(map[string]asymmetric.IPubKey) error
}

type IConfig interface {
	settings.IConfigSettingsHLS

	GetNetwork() string
	GetLogging() logger.ILogging
	GetAddress() IAddress
	GetConnections() []string
	GetFriends() map[string]asymmetric.IPubKey
	GetService(string) (string, bool)
}

type IAddress interface {
	GetTCP() string
	GetHTTP() string
}
