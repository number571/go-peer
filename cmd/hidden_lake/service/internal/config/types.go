package config

import (
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

type IWrapper interface {
	Config() IConfig
	Editor() IEditor
}

type IEditor interface {
	UpdateConnections([]string) error
	UpdateFriends(map[string]asymmetric.IPubKey) error
}

type IConfig interface {
	Network() string
	Logging() ILogging
	Address() IAddress
	Connections() []string
	Friends() map[string]asymmetric.IPubKey
	Service(string) (string, bool)
}

type ILogging interface {
	Info() bool
	Warn() bool
	Erro() bool
}

type IAddress interface {
	TCP() string
	HTTP() string
}
