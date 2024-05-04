package config

import (
	"os"

	logger "github.com/number571/go-peer/internal/logger/std"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/utils"
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
	FAddress    string       `yaml:"address"`
	FConnection *SConnection `yaml:"connection"`

	fLogging logger.ILogging
}

type SConnection struct {
	FHLTHost string `yaml:"hlt_host"`
	FPostID  string `yaml:"post_id"`
}

func BuildConfig(pFilepath string, pCfg *SConfig) (IConfig, error) {
	if _, err := os.Stat(pFilepath); !os.IsNotExist(err) {
		return nil, utils.MergeErrors(ErrConfigAlreadyExist, err)
	}

	if err := pCfg.initConfig(); err != nil {
		return nil, utils.MergeErrors(ErrInitConfig, err)
	}

	if err := os.WriteFile(pFilepath, encoding.SerializeYAML(pCfg), 0o600); err != nil {
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
		return utils.MergeErrors(ErrInvalidLogging, err)
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

func (p *SConfig) GetAddress() string {
	return p.FAddress
}
