package config

import (
	"github.com/number571/go-peer/modules/crypto/asymmetric"
)

type IConfig interface {
	Address() iAddress
	Connections() []string
	Friends() []asymmetric.IPubKey
	Service(string) (string, bool)
}

type iAddress interface {
	TCP() string
	HTTP() string
}
