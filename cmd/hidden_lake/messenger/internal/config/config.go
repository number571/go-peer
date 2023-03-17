package config

import (
	"fmt"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/filesystem"
	"github.com/number571/go-peer/pkg/logger"
)

var (
	_ IConfig     = &SConfig{}
	_ IAddress    = &SAddress{}
	_ IConnection = &SConnection{}
)

type SConfig struct {
	FLogging    []string     `json:"logging,omitempty"`
	FAddress    *SAddress    `json:"address"`
	FConnection *SConnection `json:"connection"`
	FStorageKey string       `json:"storage_key,omitempty"`

	fLogging *sLogging
}

type sLogging []bool

type SAddress struct {
	FInterface string `json:"interface"`
	FIncoming  string `json:"incoming"`
}

type SConnection struct {
	FService string `json:"service"`
	FTraffic string `json:"traffic,omitempty"`
}

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

func (cfg *SConfig) GetAddress() IAddress {
	return cfg.FAddress
}

func (cfg *SConfig) GetConnection() IConnection {
	return cfg.FConnection
}

func (cfg *SConfig) GetStorageKey() string {
	return cfg.FStorageKey
}

func (conn *SConnection) GetService() string {
	return conn.FService
}

func (conn *SConnection) GetTraffic() string {
	return conn.FTraffic
}

func (addr *SAddress) GetInterface() string {
	return addr.FInterface
}

func (addr *SAddress) GetIncoming() string {
	return addr.FIncoming
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
