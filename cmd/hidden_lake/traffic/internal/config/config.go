package config

import (
	"fmt"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/filesystem"
)

const (
	CLogInfo = "info"
	CLogWarn = "warn"
	CLogErro = "erro"
)

var (
	_ IConfig = &SConfig{}
)

type SConfig struct {
	FLogging    []string `json:"logging,omitempty"`
	FNetwork    string   `json:"network,omitempty"`
	FAddress    string   `json:"address,omitempty"`
	FConnection string   `json:"connection,omitempty"`

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

func (cfg *SConfig) Network() string {
	return cfg.FNetwork
}

func (cfg *SConfig) Address() string {
	return cfg.FAddress
}

func (cfg *SConfig) Connection() string {
	return cfg.FConnection
}

func (cfg *SConfig) Logging() ILogging {
	return cfg.fLogging
}

func (logging *sLogging) Info() bool {
	return (*logging)[0]
}

func (logging *sLogging) Warn() bool {
	return (*logging)[1]
}

func (logging *sLogging) Erro() bool {
	return (*logging)[2]
}
