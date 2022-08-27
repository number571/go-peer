package main

import (
	"github.com/number571/go-peer/cmd/hmc/config"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
)

var (
	_ iWrapper       = &sWrapper{}
	_ iConfigWrapper = &sConfigWrapper{}
)

// Wrapper

type sWrapper struct {
	cfg iConfigWrapper
}

func newWrapper(cfg iConfigWrapper) iWrapper {
	return &sWrapper{
		cfg: cfg,
	}
}

func (wrap *sWrapper) Config() iConfigWrapper {
	return wrap.cfg
}

// Config wrapper

type sConfigWrapper struct {
	cfg config.IConfig
}

func newConfigWrapper(cfg config.IConfig) iConfigWrapper {
	return &sConfigWrapper{
		cfg: cfg,
	}
}

func (cfgw *sConfigWrapper) Original() config.IConfig {
	return cfgw.cfg
}

func (cfgw *sConfigWrapper) GetNameByPubKey(pubKey asymmetric.IPubKey) (string, bool) {
	for _, friend := range cfgw.cfg.Friends() {
		if friend.PubKey().Address().String() == pubKey.Address().String() {
			return friend.Name(), true
		}
	}
	return "", false
}

func (cfgw *sConfigWrapper) GetPubKeyByName(name string) (asymmetric.IPubKey, bool) {
	for _, friend := range cfgw.cfg.Friends() {
		if friend.Name() == name {
			return friend.PubKey(), true
		}
	}
	return nil, false
}
