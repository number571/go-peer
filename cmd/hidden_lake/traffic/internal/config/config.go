package config

import (
	"fmt"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/filesystem"
	"github.com/number571/go-peer/pkg/logger"
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

func BuildConfig(pFilepath string, pCfg *SConfig) (IConfig, error) {
	configFile := filesystem.OpenFile(pFilepath)

	if configFile.IsExist() {
		return nil, errors.NewError(fmt.Sprintf("config file '%s' already exist", pFilepath))
	}

	if err := configFile.Write(encoding.Serialize(pCfg)); err != nil {
		return nil, errors.WrapError(err, "write config")
	}

	if err := pCfg.loadLogging(); err != nil {
		return nil, errors.WrapError(err, "load logging")
	}
	return pCfg, nil
}

func LoadConfig(pFilepath string) (IConfig, error) {
	configFile := filesystem.OpenFile(pFilepath)

	if !configFile.IsExist() {
		return nil, errors.NewError(fmt.Sprintf("config file '%s' does not exist", pFilepath))
	}

	bytes, err := configFile.Read()
	if err != nil {
		return nil, errors.WrapError(err, "read config")
	}

	cfg := new(SConfig)
	if err := encoding.Deserialize(bytes, cfg); err != nil {
		return nil, errors.WrapError(err, "deserialize config")
	}

	if err := cfg.loadLogging(); err != nil {
		return nil, errors.WrapError(err, "load logging")
	}
	return cfg, nil
}

func (p *SConfig) loadLogging() error {
	// [info, warn, erro]
	logging := sLogging(make([]bool, 3))

	mapping := map[string]int{
		"info": 0,
		"warn": 1,
		"erro": 2,
	}

	for _, v := range p.FLogging {
		logType, ok := mapping[v]
		if !ok {
			return errors.NewError(fmt.Sprintf("undefined log type '%s'", v))
		}
		logging[logType] = true
	}

	p.fLogging = &logging
	return nil
}

func (p *SConfig) GetNetwork() string {
	return p.FNetwork
}

func (p *SConfig) GetAddress() string {
	return p.FAddress
}

func (p *SConfig) GetConnection() string {
	return p.FConnection
}

func (p *SConfig) GetConsumers() []string {
	return p.FConsumers
}

func (p *SConfig) GetLogging() logger.ILogging {
	return p.fLogging
}

func (p *sLogging) HasInfo() bool {
	return (*p)[0]
}

func (p *sLogging) HasWarn() bool {
	return (*p)[1]
}

func (p *sLogging) HasErro() bool {
	return (*p)[2]
}
