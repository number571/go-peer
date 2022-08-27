package main

import (
	"github.com/number571/go-peer/cmd/hmc/config"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
)

type iActions map[string]iAction

type iAction interface {
	Description() string
	Do()
}

type iWrapper interface {
	Config() iConfigWrapper
}

type iConfigWrapper interface {
	Original() config.IConfig
	GetNameByPubKey(asymmetric.IPubKey) (string, bool)
	GetPubKeyByName(string) (asymmetric.IPubKey, bool)
}
