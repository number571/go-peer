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
	Logging() bool
	Network() string
	Address() iAddress
	Connections() []string
	Friends() map[string]asymmetric.IPubKey
	Service(string) (string, bool)
}

type iAddress interface {
	TCP() string
	HTTP() string
}
