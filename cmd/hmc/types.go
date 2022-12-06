package main

import (
	"github.com/number571/go-peer/cmd/hmc/config"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
)

type iWrapper interface {
	Config() iConfigWrapper
}

type iConfigWrapper interface {
	Original() config.IConfig
	GetNameByPubKey(asymmetric.IPubKey) (string, bool)
	GetPubKeyByName(string) (asymmetric.IPubKey, bool)
}

type iInputter interface {
	String() string
	Password() string
}
