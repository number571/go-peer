package config

import (
	"github.com/number571/go-peer/crypto/asymmetric"
)

type IConfig interface {
	F2F() bool
	Connections() []string
	Friends() []iFriend
}

type iFriend interface {
	Name() string
	PubKey() asymmetric.IPubKey
}
