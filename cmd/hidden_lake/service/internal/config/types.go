package config

import (
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
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
