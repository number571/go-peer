package config

import (
	"fmt"
	"sync"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/utils"
	logger "github.com/number571/go-peer/internal/logger/std"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/filesystem"
)

var (
	_ IConfigSettings = &SConfigSettings{}
	_ IConfig         = &SConfig{}
	_ IAddress        = &SAddress{}
)

type SConfigSettings struct {
	FMessagesCapacity uint64 `json:"messages_capacity"`
}

type SConfig struct {
	FSettings *SConfigSettings `json:"settings"`

	FLogging           []string  `json:"logging,omitempty"`
	FLanguage          string    `json:"language,omitempty"`
	FAddress           *SAddress `json:"address"`
	FConnection        string    `json:"connection"`
	FStorageKey        string    `json:"storage_key,omitempty"`
	FBackupConnections []string  `json:"backup_connections,omitempty"`

	fFilepath string
	fMutex    sync.Mutex
	fLanguage utils.ILanguage
	fLogging  *sLogging
}

type sLogging []bool

type SAddress struct {
	FInterface string `json:"interface"`
	FIncoming  string `json:"incoming"`
	FPPROF     string `json:"pprof,omitempty"`
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
		return nil, errors.WrapError(err, "internal init config")
	}

	return cfg, nil
}

func (p *SConfig) GetSettings() IConfigSettings {
	return p.FSettings
}

func (p *SConfigSettings) GetMessagesCapacity() uint64 {
	return p.FMessagesCapacity
}

func (p *SConfig) isValid() bool {
	return true &&
		p.FSettings.FMessagesCapacity != 0
}

func (p *SConfig) initConfig() error {
	if !p.isValid() {
		return errors.NewError("load config settings")
	}

	if err := p.loadLogging(); err != nil {
		return errors.WrapError(err, "load logging")
	}

	if err := p.loadLanguage(); err != nil {
		return errors.WrapError(err, "load language")
	}

	if p.FAddress == nil {
		p.FAddress = new(SAddress)
	}

	if p.FSettings == nil {
		p.FSettings = new(SConfigSettings)
	}

	return nil
}

func (p *SConfig) loadLanguage() error {
	res, err := utils.ToILanguage(p.FLanguage)
	if err != nil {
		return err
	}
	p.fLanguage = res
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

func (p *SConfig) GetLanguage() utils.ILanguage {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return p.fLanguage
}

func (p *SConfig) GetBackupConnections() []string {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return p.FBackupConnections
}

func (p *SConfig) GetAddress() IAddress {
	return p.FAddress
}

func (p *SConfig) GetConnection() string {
	return p.FConnection
}

func (p *SConfig) GetStorageKey() string {
	return p.FStorageKey
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

func (p *sLogging) HasInfo() bool {
	return (*p)[0]
}

func (p *sLogging) HasWarn() bool {
	return (*p)[1]
}

func (p *sLogging) HasErro() bool {
	return (*p)[2]
}
