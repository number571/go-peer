package config

import (
	"errors"
	"fmt"
	"os"
	"sync"

	logger "github.com/number571/go-peer/internal/logger/std"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
)

var (
	_ IConfigSettings = &SConfigSettings{}
	_ IConfig         = &SConfig{}
	_ IAddress        = &SAddress{}
	_ logger.ILogging = &sLogging{}
)

type SConfigSettings struct {
	fMutex              sync.Mutex
	FMessageSizeBytes   uint64 `yaml:"message_size_bytes"`
	FWorkSizeBits       uint64 `yaml:"work_size_bits"`
	FQueuePeriodMS      uint64 `yaml:"queue_period_ms"`
	FKeySizeBits        uint64 `yaml:"key_size_bits"`
	FLimitVoidSizeBytes uint64 `yaml:"limit_void_size_bytes,omitempty"`
	FNetworkKey         string `yaml:"network_key,omitempty"`
}

type SConfig struct {
	FSettings *SConfigSettings `yaml:"settings"`

	FLogging     []string             `yaml:"logging,omitempty"`
	FAddress     *SAddress            `yaml:"address,omitempty"`
	FConnections []string             `yaml:"connections,omitempty"`
	FServices    map[string]*SService `yaml:"services,omitempty"`
	FFriends     map[string]string    `yaml:"friends,omitempty"`

	fFilepath string
	fMutex    sync.Mutex
	fLogging  *sLogging
	fFriends  map[string]asymmetric.IPubKey
}

type sLogging []bool

type SService struct {
	FHost  string `yaml:"host,omitempty"`
	FShare bool   `yaml:"share,omitempty"`
}

type SAddress struct {
	FTCP   string `yaml:"tcp,omitempty"`
	FHTTP  string `yaml:"http,omitempty"`
	FPPROF string `yaml:"pprof,omitempty"`
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
		return nil, fmt.Errorf("load logging: %w", err)
	}

	return cfg, nil
}

func (p *SConfigSettings) GetMessageSizeBytes() uint64 {
	return p.FMessageSizeBytes
}

func (p *SConfigSettings) GetWorkSizeBits() uint64 {
	return p.FWorkSizeBits
}

func (p *SConfigSettings) GetKeySizeBits() uint64 {
	return p.FKeySizeBits
}

func (p *SConfigSettings) GetQueuePeriodMS() uint64 {
	return p.FQueuePeriodMS
}

func (p *SConfigSettings) GetLimitVoidSizeBytes() uint64 {
	return p.FLimitVoidSizeBytes
}

func (p *SConfigSettings) GetNetworkKey() string {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return p.FNetworkKey
}

func (p *SConfig) GetSettings() IConfigSettings {
	return p.FSettings
}

func (p *SConfig) isValid() bool {
	for _, v := range p.FServices {
		if v.FHost == "" && !v.FShare {
			return false
		}
	}
	return true &&
		p.FSettings.FMessageSizeBytes != 0 &&
		p.FSettings.FKeySizeBits != 0 &&
		p.FSettings.FQueuePeriodMS != 0
}

func (p *SConfig) initConfig() error {
	if p.FSettings == nil {
		p.FSettings = new(SConfigSettings)
	}

	if p.FAddress == nil {
		p.FAddress = new(SAddress)
	}

	for k, v := range p.FServices {
		if v == nil {
			p.FServices[k] = new(SService)
		}
	}

	if !p.isValid() {
		return errors.New("load config settings")
	}

	if err := p.loadPubKeys(); err != nil {
		return fmt.Errorf("load public keys: %w", err)
	}

	if err := p.loadLogging(); err != nil {
		return fmt.Errorf("load logging: %w", err)
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
			return fmt.Errorf("undefined log type '%s'", v)
		}
		logging[logType] = true
	}

	p.fLogging = &logging
	return nil
}

func (p *SConfig) loadPubKeys() error {
	p.fFriends = make(map[string]asymmetric.IPubKey)

	mapping := make(map[string]struct{})
	for name, val := range p.FFriends {
		if _, ok := mapping[val]; ok {
			return fmt.Errorf("found public key duplicate '%s'", val)
		}
		mapping[val] = struct{}{}

		pubKey := asymmetric.LoadRSAPubKey(val)
		if pubKey == nil {
			return fmt.Errorf("public key is nil for '%s'", name)
		}

		p.fFriends[name] = pubKey
		if pubKey.GetSize() != p.FSettings.FKeySizeBits {
			return fmt.Errorf("not supported key size for '%s'", name)
		}
	}

	return nil
}

func (p *SConfig) GetFriends() map[string]asymmetric.IPubKey {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	result := make(map[string]asymmetric.IPubKey)
	for k, v := range p.fFriends {
		result[k] = v
	}
	return result
}

func (p *SConfig) GetLogging() logger.ILogging {
	return p.fLogging
}

func (p *SConfig) GetAddress() IAddress {
	return p.FAddress
}

func (p *SConfig) GetConnections() []string {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return p.FConnections
}

func (p *SConfig) GetService(name string) (IService, bool) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	service, ok := p.FServices[name]
	return service, ok
}

func (p *SService) GetHost() string {
	return p.FHost
}

func (p *SService) GetShare() bool {
	return p.FShare
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

func (p *SAddress) GetTCP() string {
	return p.FTCP
}

func (p *SAddress) GetHTTP() string {
	return p.FHTTP
}

func (p *SAddress) GetPPROF() string {
	return p.FPPROF
}
