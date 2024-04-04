package config

import (
	"os"

	logger "github.com/number571/go-peer/internal/logger/std"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/utils"
)

var (
	_ IConfigSettings = &SConfigSettings{}
	_ IConfig         = &SConfig{}
	_ IAddress        = &SAddress{}
)

type SConfigSettings struct {
	FMessageSizeBytes   uint64 `json:"message_size_bytes" yaml:"message_size_bytes"`
	FKeySizeBits        uint64 `json:"key_size_bits" yaml:"key_size_bits"`
	FMessagesCapacity   uint64 `json:"messages_capacity" yaml:"messages_capacity"`
	FWorkSizeBits       uint64 `json:"work_size_bits,omitempty" yaml:"work_size_bits,omitempty"`
	FQueuePeriodMS      uint64 `json:"queue_period_ms,omitempty" yaml:"queue_period_ms,omitempty"`
	FLimitVoidSizeBytes uint64 `json:"limit_void_size_bytes,omitempty" yaml:"limit_void_size_bytes,omitempty"`
	FNetworkKey         string `json:"network_key,omitempty" yaml:"network_key,omitempty"`
	FStorageEnabled     bool   `json:"storage_enabled,omitempty" yaml:"storage_enabled,omitempty"`
}

type SConfig struct {
	FSettings *SConfigSettings `yaml:"settings"`

	FLogging     []string  `yaml:"logging,omitempty"`
	FAddress     *SAddress `yaml:"address,omitempty"`
	FConnections []string  `yaml:"connections,omitempty"`
	FConsumers   []string  `yaml:"consumers,omitempty"`

	fLogging logger.ILogging
}

type SAddress struct {
	FTCP   string `yaml:"tcp,omitempty"`
	FHTTP  string `yaml:"http,omitempty"`
	FPPROF string `yaml:"pprof,omitempty"`
}

func BuildConfig(pFilepath string, pCfg *SConfig) (IConfig, error) {
	if _, err := os.Stat(pFilepath); !os.IsNotExist(err) {
		return nil, utils.MergeErrors(ErrConfigAlreadyExist, err)
	}

	if err := pCfg.initConfig(); err != nil {
		return nil, utils.MergeErrors(ErrInitConfig, err)
	}

	if err := os.WriteFile(pFilepath, encoding.SerializeYAML(pCfg), 0o644); err != nil {
		return nil, utils.MergeErrors(ErrWriteConfig, err)
	}

	return pCfg, nil
}

func LoadConfig(pFilepath string) (IConfig, error) {
	if _, err := os.Stat(pFilepath); os.IsNotExist(err) {
		return nil, utils.MergeErrors(ErrConfigNotExist, err)
	}

	bytes, err := os.ReadFile(pFilepath)
	if err != nil {
		return nil, utils.MergeErrors(ErrReadConfig, err)
	}

	cfg := new(SConfig)
	if err := encoding.DeserializeYAML(bytes, cfg); err != nil {
		return nil, utils.MergeErrors(ErrDeserializeConfig, err)
	}

	if err := cfg.initConfig(); err != nil {
		return nil, utils.MergeErrors(ErrInitConfig, err)
	}
	return cfg, nil
}

func (p *SConfig) GetSettings() IConfigSettings {
	return p.FSettings
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

func (p *SConfigSettings) GetNetworkKey() string {
	return p.FNetworkKey
}

func (p *SConfigSettings) GetKeySizeBits() uint64 {
	return p.FKeySizeBits
}

func (p *SConfigSettings) GetStorageEnabled() bool {
	return p.FStorageEnabled
}

func (p *SConfig) isValid() bool {
	return true &&
		p.FSettings.FMessageSizeBytes != 0 &&
		p.FSettings.FMessagesCapacity != 0 &&
		p.FSettings.FKeySizeBits != 0
}

func (p *SConfig) initConfig() error {
	if p.FSettings == nil {
		p.FSettings = new(SConfigSettings)
	}

	if p.FAddress == nil {
		p.FAddress = new(SAddress)
	}

	if !p.isValid() {
		return ErrInvalidConfig
	}

	if err := p.loadLogging(); err != nil {
		return utils.MergeErrors(ErrLoadLogging, err)
	}

	return nil
}

func (p *SConfig) loadLogging() error {
	result, err := logger.LoadLogging(p.FLogging)
	if err != nil {
		return utils.MergeErrors(ErrLoadLogging, err)
	}
	p.fLogging = result
	return nil
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
