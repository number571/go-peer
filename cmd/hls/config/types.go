package config

import (
	"github.com/number571/go-peer/modules/crypto/asymmetric"
)

type IEditor interface {
	UpdateConnections([]string) error
	UpdateFriends([]asymmetric.IPubKey) error
}

type IConfig interface {
	Network() string
	Address() iAddress
	Connections() []string
	Friends() []asymmetric.IPubKey
	Service(string) (string, bool)
}

type iAddress interface {
	TCP() string
	HTTP() string
}
