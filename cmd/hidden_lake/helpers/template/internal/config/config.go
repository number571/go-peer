package config

import (
	"os"

	logger "github.com/number571/go-peer/internal/logger/std"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/utils"
)

var (
	_ IConfig  = &SConfig{}
	_ IAddress = &SAddress{}
)

type SConfigSettings struct {
	FValue string `json:"value" yaml:"value"`
	// TODO: need implementation
}

type SConfig struct {
	FSettings *SConfigSettings `yaml:"settings"`

	FLogging []string  `yaml:"logging,omitempty"`
	FAddress *SAddress `yaml:"address"`

	fLogging logger.ILogging
}

type SAddress struct {
	FHTTP  string `yaml:"http"`
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

func (p *SConfig) isValid() bool {
	return true &&
		p.FAddress.FHTTP != ""
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

func (p *SConfig) GetSettings() IConfigSettings {
	return p.FSettings
}

func (p *SConfigSettings) GetValue() string {
	return p.FValue
}

func (p *SConfig) GetAddress() IAddress {
	return p.FAddress
}

func (p *SAddress) GetHTTP() string {
	return p.FHTTP
}

func (p *SAddress) GetPPROF() string {
	return p.FPPROF
}

func (p *SConfig) GetLogging() logger.ILogging {
	return p.fLogging
}
