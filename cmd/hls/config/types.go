package config

import "github.com/number571/go-peer/crypto"

type IConfig interface {
	CleanCron() string
	Address() iAddress
	F2F() iF2F
	Connections() []string
	CheckOnline() []crypto.IPubKey
	GetService(string) (iBlock, bool)
}

type iF2F interface {
	Status() bool
	Friends() []crypto.IPubKey
}

type iAddress interface {
	HLS() string
	HTTP() string
}

type iBlock interface {
	Address() string
	IsRedirect() bool
}
