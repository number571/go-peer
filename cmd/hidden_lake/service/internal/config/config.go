package config

import (
	"fmt"
	"sync"

	logger "github.com/number571/go-peer/internal/logger/std"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/filesystem"
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
	FQueuePeriodMS      uint64 `json:"queue_period_ms"`
	FKeySizeBits        uint64 `json:"key_size_bits"`
	FLimitVoidSizeBytes uint64 `json:"limit_void_size_bytes,omitempty"`
}

type SConfig struct {
	FSettings *SConfigSettings `json:"settings"`

	FLogging     []string          `json:"logging,omitempty"`
	FAddress     *SAddress         `json:"address,omitempty"`
	FNetworkKey  string            `json:"network_key,omitempty"`
	FConnections []string          `json:"connections,omitempty"`
	FServices    map[string]string `json:"services,omitempty"`
	FFriends     map[string]string `json:"friends,omitempty"`

	fFilepath string
	fMutex    sync.Mutex
	fLogging  *sLogging
	fFriends  map[string]asymmetric.IPubKey
}

type sLogging []bool

type SAddress struct {
	FTCP   string `json:"tcp,omitempty"`
	FHTTP  string `json:"http,omitempty"`
	FPPROF string `json:"pprof,omitempty"`
}

func BuildConfig(pFilepath string, pCfg *SConfig) (IConfig, error) {
	configFile := filesystem.OpenFile(pFilepath)
	if configFile.IsExist() {
		return nil, errors.NewError(fmt.Sprintf("config file '%s' already exist", pFilepath))
	}

	if err := configFile.Write(encoding.Serialize(pCfg, true)); err != nil {
		return nil, errors.WrapError(err, "write config")
	}

	pCfg.fFilepath = pFilepath
	if err := pCfg.initConfig(); err != nil {
		return nil, errors.WrapError(err, "init config")
	}
	return pCfg, nil
}

func LoadConfig(pFilepath string) (IConfig, error) {
	configFile := filesystem.OpenFile(pFilepath)

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

	cfg.fFilepath = pFilepath
	if err := cfg.initConfig(); err != nil {
		return nil, errors.WrapError(err, "init config")
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

func (p *SConfig) GetSettings() IConfigSettings {
	return p.FSettings
}

func (p *SConfig) isValid() bool {
	return true &&
		p.FSettings.FMessageSizeBytes != 0 &&
		p.FSettings.FWorkSizeBits != 0 &&
		p.FSettings.FKeySizeBits != 0 &&
		p.FSettings.FQueuePeriodMS != 0
}

func (p *SConfig) initConfig() error {
	if !p.isValid() {
		return errors.NewError("load config settings")
	}

	if err := p.loadPubKeys(); err != nil {
		return errors.WrapError(err, "load public keys")
	}

	if err := p.loadLogging(); err != nil {
		return errors.WrapError(err, "load logging")
	}

	if p.FSettings == nil {
		p.FSettings = new(SConfigSettings)
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

func (p *SConfig) loadPubKeys() error {
	p.fFriends = make(map[string]asymmetric.IPubKey)

	mapping := make(map[string]interface{})
	for name, val := range p.FFriends {
		if _, ok := mapping[val]; ok {
			return fmt.Errorf("found public key duplicate '%s'", val)
		}
		pubKey := asymmetric.LoadRSAPubKey(val)
		if pubKey == nil {
			return errors.NewError(fmt.Sprintf("public key is nil for '%s'", name))
		}
		p.fFriends[name] = pubKey
		if pubKey.GetSize() != p.FSettings.FKeySizeBits {
			return errors.NewError(fmt.Sprintf("not supported key size for '%s'", name))
		}
	}

	return nil
}

func (p *SConfig) GetNetworkKey() string {
	return p.FNetworkKey
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
	if p == nil {
		return &SAddress{}
	}
	return p.FAddress
}

func (p *SConfig) GetConnections() []string {
	return p.FConnections
}

func (p *SConfig) GetService(name string) (string, bool) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	addr, ok := p.FServices[name]
	return addr, ok
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
