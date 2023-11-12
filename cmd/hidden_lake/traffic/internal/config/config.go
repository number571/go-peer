package config

import (
	"fmt"

	logger "github.com/number571/go-peer/internal/logger/std"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/file_system"
)

var (
	_ IConfigSettings = &SConfigSettings{}
	_ IConfig         = &SConfig{}
	_ IAddress        = &SAddress{}
	_ logger.ILogging = &sLogging{}
)

type SConfigSettings struct {
	FMessageSizeBytes   uint64 `json:"message_size_bytes"`
	FWorkSizeBits       uint64 `json:"work_size_bits"`
	FMessagesCapacity   uint64 `json:"messages_capacity"`
	FQueuePeriodMS      uint64 `json:"queue_period_ms,omitempty"`
	FLimitVoidSizeBytes uint64 `json:"limit_void_size_bytes,omitempty"`
}

type SConfig struct {
	FSettings *SConfigSettings `json:"settings"`

	FLogging     []string  `json:"logging,omitempty"`
	FAddress     *SAddress `json:"address,omitempty"`
	FIsStorage   bool      `json:"is_storage,omitempty"`
	FNetworkKey  string    `json:"network_key,omitempty"`
	FConnections []string  `json:"connections,omitempty"`
	FConsumers   []string  `json:"consumers,omitempty"`

	fLogging *sLogging
}

type SAddress struct {
	FTCP   string `json:"tcp,omitempty"`
	FHTTP  string `json:"http,omitempty"`
	FPPROF string `json:"pprof,omitempty"`
}

type sLogging []bool

func BuildConfig(pFilepath string, pCfg *SConfig) (IConfig, error) {
	configFile := file_system.OpenFile(pFilepath)

	if configFile.IsExist() {
		return nil, errors.NewError(fmt.Sprintf("config file '%s' already exist", pFilepath))
	}

	if err := configFile.Write(encoding.Serialize(pCfg, true)); err != nil {
		return nil, errors.WrapError(err, "write config")
	}

	if err := pCfg.initConfig(); err != nil {
		return nil, errors.WrapError(err, "load logging")
	}
	return pCfg, nil
}

func LoadConfig(pFilepath string) (IConfig, error) {
	configFile := file_system.OpenFile(pFilepath)

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

	if err := cfg.initConfig(); err != nil {
		return nil, errors.WrapError(err, "load logging")
	}
	return cfg, nil
}

func (p *SConfig) GetSettings() IConfigSettings {
	return p.FSettings
}

func (p *SConfig) GetIsStorage() bool {
	return p.FIsStorage
}

func (p *SConfigSettings) GetMessageSizeBytes() uint64 {
	return p.FMessageSizeBytes
}

func (p *SConfigSettings) GetWorkSizeBits() uint64 {
	return p.FWorkSizeBits
}

func (p *SConfigSettings) GetQueuePeriodMS() uint64 {
	return p.FQueuePeriodMS
}

func (p *SConfigSettings) GetMessagesCapacity() uint64 {
	return p.FMessagesCapacity
}

func (p *SConfigSettings) GetLimitVoidSizeBytes() uint64 {
	return p.FLimitVoidSizeBytes
}

func (p *SConfig) isValid() bool {
	return true &&
		p.FSettings.FMessageSizeBytes != 0 &&
		p.FSettings.FMessagesCapacity != 0
}

func (p *SConfig) initConfig() error {
	if p.FSettings == nil {
		p.FSettings = new(SConfigSettings)
	}

	if p.FAddress == nil {
		p.FAddress = new(SAddress)
	}

	if !p.isValid() {
		return errors.NewError("load config settings")
	}

	if err := p.loadLogging(); err != nil {
		return errors.WrapError(err, "load logging")
	}

	return nil
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

func (p *SConfig) GetNetworkKey() string {
	return p.FNetworkKey
}

func (p *SConfig) GetAddress() IAddress {
	return p.FAddress
}

func (p *SAddress) GetTCP() string {
	return p.FTCP
}

func (p *SAddress) GetHTTP() string {
	return p.FHTTP
}

func (p *SAddress) GetPPROF() string {
	return p.FPPROF
}

func (p *SConfig) GetConnections() []string {
	return p.FConnections
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
