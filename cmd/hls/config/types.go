package config

import (
	"github.com/number571/go-peer/crypto/asymmetric"
)

type IConfig interface {
	CleanCron() string
	Address() iAddress
	Connections() []string
	F2F() iF2F
	OnlineChecker() iOnlineChecker
	GetService(string) (iBlock, bool)
}

type iOnlineChecker interface {
	Status() bool
	PubKeys() []asymmetric.IPubKey
}

type iF2F interface {
	Status() bool
	PubKeys() []asymmetric.IPubKey
}

type iAddress interface {
	HLS() string
	HTTP() string
}

type iBlock interface {
	Address() string
	IsRedirect() bool
}
