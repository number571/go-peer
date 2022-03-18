package config

import (
	"sync"

	"github.com/number571/go-peer/cmd/hls/utils"
)

var (
	_ IConfig = &sConfig{}
)

type sConfig struct {
	fMutex    sync.Mutex
	FAddress  string            `json:"address"`
	FConnects []string          `json:"connects"`
	FServices map[string]string `json:"services"`
}

const (
	// create local hls
	cAddress = "localhost:9571"
)

var (
	// connect to another hls's
	gConnects = []string{"127.0.0.2:9571"}

	// crypto-address -> network-address
	gServices = map[string]string{
		"hidden-default-service": "http://localhost:8080",
	}
)

func NewConfig(filepath string) IConfig {
	var cfg = new(sConfig)

	if !utils.FileIsExist(filepath) {
		cfg = &sConfig{
			FAddress:  cAddress,
			FConnects: gConnects,
			FServices: gServices,
		}
		err := utils.WriteFile(filepath, utils.Serialize(cfg))
		if err != nil {
			panic(err)
		}
	} else {
		err := utils.Deserialize(utils.ReadFile(filepath), cfg)
		if err != nil {
			panic(err)
		}
	}

	return cfg
}

func (cfg *sConfig) Address() string {
	return cfg.FAddress
}

func (cfg *sConfig) Connections() []string {
	return cfg.FConnects
}

func (cfg *sConfig) GetService(name string) (string, bool) {
	cfg.fMutex.Lock()
	defer cfg.fMutex.Unlock()

	addr, ok := cfg.FServices[name]
	return addr, ok
}
