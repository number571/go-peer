package config

import (
	"github.com/number571/go-peer/modules/encoding"
	"github.com/number571/go-peer/modules/filesystem"
)

var (
	_ IConfig = &sConfig{}
)

type sConfig struct {
	FAddress    string `json:"address"`
	FConnection string `json:"connection"`
}

const (
	// create local hlm
	cDefaultAddress = "localhost:8080"

	// create local hls
	cDefaultConnection = "localhost:9572"
)

func NewConfig(filepath string) IConfig {
	var cfg = new(sConfig)

	if !filesystem.OpenFile(filepath).IsExist() {
		cfg = &sConfig{
			FAddress:    cDefaultAddress,
			FConnection: cDefaultConnection,
		}
		err := filesystem.OpenFile(filepath).Write(encoding.Serialize(cfg))
		if err != nil {
			panic(err)
		}
	} else {
		bytes, err := filesystem.OpenFile(filepath).Read()
		if err != nil {
			panic(err)
		}
		err = encoding.Deserialize(bytes, cfg)
		if err != nil {
			panic(err)
		}
	}

	return cfg
}

func (cfg *sConfig) Address() string {
	return cfg.FAddress
}

func (cfg *sConfig) Connection() string {
	return cfg.FConnection
}
