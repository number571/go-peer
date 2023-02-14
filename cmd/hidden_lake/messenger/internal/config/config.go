package config

import (
	"fmt"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/filesystem"
)

var (
	_ IConfig     = &SConfig{}
	_ iAddress    = &SAddress{}
	_ iConnection = &SConnection{}
)

type SConfig struct {
	FAddress    *SAddress    `json:"address"`
	FConnection *SConnection `json:"connection"`
	FStorageKey string       `json:"storage_key,omitempty"`
}

type SAddress struct {
	FInterface string `json:"interface"`
	FIncoming  string `json:"incoming"`
}

type SConnection struct {
	FService string `json:"service"`
	FTraffic string `json:"traffic,omitempty"`
}

func NewConfig(filepath string, cfg *SConfig) (IConfig, error) {
	configFile := filesystem.OpenFile(filepath)

	if configFile.IsExist() {
		return nil, fmt.Errorf("config file '%s' already exist", filepath)
	}

	if err := configFile.Write(encoding.Serialize(cfg)); err != nil {
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

	return cfg, nil
}

func (cfg *SConfig) Address() iAddress {
	return cfg.FAddress
}

func (cfg *SConfig) Connection() iConnection {
	return cfg.FConnection
}

func (cfg *SConfig) StorageKey() string {
	return cfg.FStorageKey
}

func (conn *SConnection) Service() string {
	return conn.FService
}

func (conn *SConnection) Traffic() string {
	return conn.FTraffic
}

func (addr *SAddress) Interface() string {
	return addr.FInterface
}

func (addr *SAddress) Incoming() string {
	return addr.FIncoming
}
