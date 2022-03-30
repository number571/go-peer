package config

import "github.com/number571/go-peer/crypto"

type IConfig interface {
	F2F() bool
	Address() string
	Connections() []string
	PubKeys() []crypto.IPubKey
	GetService(string) (iBlock, bool)
	CleanCron() string
}

type iBlock interface {
	Address() string
	IsRedirect() bool
}
