package config

import (
	"fmt"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/filesystem"
	"github.com/number571/go-peer/pkg/logger"
)

const (
	CLogInfo = "info"
	CLogWarn = "warn"
	CLogErro = "erro"
)

var (
	_ IConfig         = &SConfig{}
	_ logger.ILogging = &sLogging{}
)

type SConfig struct {
	FLogging    []string `json:"logging,omitempty"`
	FNetwork    string   `json:"network,omitempty"`
	FAddress    string   `json:"address,omitempty"`
	FConnection string   `json:"connection,omitempty"`
	FConsumers  []string `json:"consumers,omitempty"`

	fLogging *sLogging
}

type sLogging []bool

func BuildConfig(filepath string, cfg *SConfig) (IConfig, error) {
	configFile := filesystem.OpenFile(filepath)

	if configFile.IsExist() {
		return nil, fmt.Errorf("config file '%s' already exist", filepath)
	}

	if err := configFile.Write(encoding.Serialize(cfg)); err != nil {
		return nil, err
	}

	if err := cfg.loadLogging(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func LoadConfig(filepath string) (IConfig, error) {
	configFile := filesystem.OpenFile(filepath)

	if !configFile.IsExist() {
		return nil, fmt.Errorf("config file '%s' does not exist", filepath)
	}

	bytes, err := configFile.Read()
	if err != nil {
		return nil, err
	}

	cfg := new(SConfig)
	if err := encoding.Deserialize(bytes, cfg); err != nil {
		return nil, err
	}

	if err := cfg.loadLogging(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (cfg *SConfig) loadLogging() error {
	// [info, warn, erro]
	logging := sLogging(make([]bool, 3))

	mapping := map[string]int{
		"info": 0,
		"warn": 1,
		"erro": 2,
	}

	for _, v := range cfg.FLogging {
		logType, ok := mapping[v]
		if !ok {
			return fmt.Errorf("undefined log type '%s'", v)
		}
		logging[logType] = true
	}

	cfg.fLogging = &logging
	return nil
}

func (cfg *SConfig) GetNetwork() string {
	return cfg.FNetwork
}

func (cfg *SConfig) GetAddress() string {
	return cfg.FAddress
}

func (cfg *SConfig) GetConnection() string {
	return cfg.FConnection
}

func (cfg *SConfig) GetConsumers() []string {
	return cfg.FConsumers
}

func (cfg *SConfig) GetLogging() logger.ILogging {
	return cfg.fLogging
}

func (logging *sLogging) HasInfo() bool {
	return (*logging)[0]
}

func (logging *sLogging) HasWarn() bool {
	return (*logging)[1]
}

func (logging *sLogging) HasErro() bool {
	return (*logging)[2]
}
