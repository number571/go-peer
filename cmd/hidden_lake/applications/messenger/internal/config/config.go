package config

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/utils"
	"github.com/number571/go-peer/internal/language"
	logger "github.com/number571/go-peer/internal/logger/std"
	"github.com/number571/go-peer/pkg/encoding"
)

var (
	_ IConfigSettings = &SConfigSettings{}
	_ IConfig         = &SConfig{}
	_ IAddress        = &SAddress{}
)

type SConfigSettings struct {
	FMessagesCapacity uint64 `yaml:"messages_capacity"`
	FWorkSizeBits     uint64 `yaml:"work_size_bits,omitempty"`
	FPseudonym        string `yaml:"pseudonym"`
	FStorageKey       string `yaml:"storage_key,omitempty"`
	FLanguage         string `yaml:"language,omitempty"`

	fMutex    sync.RWMutex
	fLanguage language.ILanguage
}

type SConfig struct {
	FSettings *SConfigSettings `yaml:"settings"`

	FLogging    []string  `yaml:"logging,omitempty"`
	FAddress    *SAddress `yaml:"address"`
	FConnection string    `yaml:"connection"`

	fFilepath string
	fLogging  logger.ILogging
}

type SAddress struct {
	FInterface string `yaml:"interface"`
	FIncoming  string `yaml:"incoming"`
	FPPROF     string `yaml:"pprof,omitempty"`
}

func BuildConfig(pFilepath string, pCfg *SConfig) (IConfig, error) {
	if _, err := os.Stat(pFilepath); !os.IsNotExist(err) {
		return nil, fmt.Errorf("config file '%s' already exist", pFilepath)
	}

	pCfg.fFilepath = pFilepath
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

	cfg.fFilepath = pFilepath
	if err := cfg.initConfig(); err != nil {
		return nil, fmt.Errorf("internal init config: %w", err)
	}

	return cfg, nil
}

func (p *SConfig) GetSettings() IConfigSettings {
	return p.FSettings
}

func (p *SConfigSettings) GetMessagesCapacity() uint64 {
	return p.FMessagesCapacity
}

func (p *SConfigSettings) GetWorkSizeBits() uint64 {
	return p.FWorkSizeBits
}

func (p *SConfigSettings) GetPseudonym() string {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	return p.FPseudonym
}

func (p *SConfigSettings) GetStorageKey() string {
	return p.FStorageKey
}

func (p *SConfigSettings) GetLanguage() language.ILanguage {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	return p.fLanguage
}

func (p *SConfigSettings) loadLanguage() error {
	res, err := language.ToILanguage(p.FLanguage)
	if err != nil {
		return err
	}
	p.fLanguage = res
	return nil
}

func (p *SConfig) isValid() bool {
	return true &&
		utils.PseudonymIsValid(p.FSettings.FPseudonym) &&
		p.FConnection != "" &&
		p.FAddress.FInterface != "" &&
		p.FAddress.FIncoming != "" &&
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
		return errors.New("load config settings")
	}

	if err := p.loadLogging(); err != nil {
		return fmt.Errorf("load logging: %w", err)
	}

	if err := p.FSettings.loadLanguage(); err != nil {
		return fmt.Errorf("load language: %w", err)
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

func (p *SConfig) GetAddress() IAddress {
	return p.FAddress
}

func (p *SConfig) GetConnection() string {
	return p.FConnection
}

func (p *SAddress) GetInterface() string {
	return p.FInterface
}

func (p *SAddress) GetIncoming() string {
	return p.FIncoming
}

func (p *SAddress) GetPPROF() string {
	return p.FPPROF
}

func (p *SConfig) GetLogging() logger.ILogging {
	return p.fLogging
}
