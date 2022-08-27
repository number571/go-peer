package config

import (
	"github.com/number571/go-peer/crypto/asymmetric"
)

type IConfig interface {
	Address() iAddress
	Connections() []string
	Friends() []asymmetric.IPubKey
	GetService(string) (string, bool)
}

type iAddress interface {
	TCP() string
	HTTP() string
}
