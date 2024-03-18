package config

import (
	"errors"
	"fmt"
	"os"

	logger "github.com/number571/go-peer/internal/logger/std"
	"github.com/number571/go-peer/pkg/encoding"
)

var (
	_ IConfig     = &SConfig{}
	_ IConnection = &SConnection{}
)

type SConfigSettings struct {
	FWorkSizeBits uint64 `json:"work_size_bits,omitempty" yaml:"work_size_bits,omitempty"`
	FNetworkKey   string `json:"network_key,omitempty" yaml:"network_key,omitempty"`
	FWaitTimeMS   uint64 `json:"wait_time_ms" yaml:"wait_time_ms"`
}

type SConfig struct {
	FSettings *SConfigSettings `yaml:"settings"`

	FLogging    []string     `yaml:"logging,omitempty"`
	FConnection *SConnection `yaml:"connection"`

	fLogging logger.ILogging
}

type SConnection struct {
	FHLTHost string `yaml:"hlt_host"`
	FPostID  string `yaml:"post_id"`
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
		p.FConnection.FHLTHost != "" &&
		p.FConnection.FPostID != "" &&
		p.FSettings.FWaitTimeMS != 0
}

func (p *SConfig) initConfig() error {
	if p.FSettings == nil {
		p.FSettings = new(SConfigSettings)
	}

	if p.FConnection == nil {
		p.FConnection = new(SConnection)
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

func (p *SConfig) GetLogging() logger.ILogging {
	return p.fLogging
}

func (p *SConfigSettings) GetWaitTimeMS() uint64 {
	return p.FWaitTimeMS
}

func (p *SConfigSettings) GetNetworkKey() string {
	return p.FNetworkKey
}

func (p *SConfigSettings) GetWorkSizeBits() uint64 {
	return p.FWorkSizeBits
}

func (p *SConfig) GetConnection() IConnection {
	return p.FConnection
}

func (p *SConnection) GetHLTHost() string {
	return p.FHLTHost
}

func (p *SConnection) GetPostID() string {
	return p.FPostID
}
