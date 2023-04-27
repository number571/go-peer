package config

import (
	"github.com/number571/go-peer/pkg/filesystem"
)

func InitConfig(cfgPath string, initCfg *SConfig) (IConfig, error) {
	if filesystem.OpenFile(cfgPath).IsExist() {
		return LoadConfig(cfgPath)
	}
	if initCfg == nil {
		initCfg = &SConfig{
			FAddress: &SAddress{
				FInterface: "localhost:9591",
				FIncoming:  "localhost:9592",
			},
			FConnection: &SConnection{
				FService: "localhost:9572",
			},
		}
	}
	return BuildConfig(cfgPath, initCfg)
}
