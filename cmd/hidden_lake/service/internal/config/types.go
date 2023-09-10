package config

import (
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/config"
	logger "github.com/number571/go-peer/internal/logger/std"
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
	GetSettings() config.IConfigSettings
	GetLogging() logger.ILogging
	GetAddress() IAddress
	GetNetworkKey() string
	GetConnections() []string
	GetFriends() map[string]asymmetric.IPubKey
	GetService(string) (string, bool)
}

type IAddress interface {
	GetTCP() string
	GetHTTP() string
}
