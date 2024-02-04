package config

import (
	"errors"
	"fmt"
	"os"

	logger "github.com/number571/go-peer/internal/logger/std"
	"github.com/number571/go-peer/pkg/encoding"
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
		return nil, fmt.Errorf("config file '%s' already exist", pFilepath)
	}

	if err := pCfg.initConfig(); err != nil {
		return nil, fmt.Errorf("init config: %w", err)
	}

	if err := os.WriteFile(pFilepath, encoding.SerializeYAML(pCfg), 0o644); err != nil {
		return nil, fmt.Errorf("write config: %w", err)
	}

	return pCfg, nil
}

func LoadConfig(pFilepath string) (IConfig, error) {
	if _, err := os.Stat(pFilepath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file '%s' does not exist", pFilepath)
	}

	bytes, err := os.ReadFile(pFilepath)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	cfg := new(SConfig)
	if err := encoding.DeserializeYAML(bytes, cfg); err != nil {
		return nil, fmt.Errorf("deserialize config: %w", err)
	}

	if err := cfg.initConfig(); err != nil {
		return nil, fmt.Errorf("internal init config: %w", err)
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
		return errors.New("load config settings")
	}

	if err := p.loadLogging(); err != nil {
		return fmt.Errorf("load logging: %w", err)
	}

	return nil
}

func (p *SConfig) loadLogging() error {
	result, err := logger.LoadLogging(p.FLogging)
	if err != nil {
		return err
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
