package config

import (
	"fmt"
	"sync"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/logger"
	"github.com/number571/go-peer/internal/settings"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/filesystem"
)

var (
	_ IConfig         = &SConfig{}
	_ IAddress        = &SAddress{}
	_ logger.ILogging = &sLogging{}
)

type SConfig struct {
	settings.SConfigSettings

	FLogging []string `json:"logging,omitempty"`

	FNetwork string    `json:"network,omitempty"`
	FAddress *SAddress `json:"address,omitempty"`

	FConnections []string          `json:"connections,omitempty"`
	FServices    map[string]string `json:"services,omitempty"`
	FFriends     map[string]string `json:"friends,omitempty"`

	fFilepath string
	fMutex    sync.Mutex
	fLogging  *sLogging
	fFriends  map[string]asymmetric.IPubKey
}

type sLogging []bool

type STraffic struct {
	FDownload []string `json:"download,omitempty"`
	FUpload   []string `json:"upload,omitempty"`
}

type SAddress struct {
	FTCP  string `json:"tcp,omitempty"`
	FHTTP string `json:"http,omitempty"`
}

func BuildConfig(pFilepath string, pCfg *SConfig) (IConfig, error) {
	configFile := filesystem.OpenFile(pFilepath)
	if configFile.IsExist() {
		return nil, errors.NewError(fmt.Sprintf("config file '%s' already exist", pFilepath))
	}

	if err := configFile.Write(encoding.Serialize(pCfg, true)); err != nil {
		return nil, errors.WrapError(err, "write config")
	}

	if err := pCfg.initConfig(pFilepath); err != nil {
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

	if err := cfg.initConfig(pFilepath); err != nil {
		return nil, errors.WrapError(err, "init config")
	}
	return cfg, nil
}

func (p *SConfig) initConfig(filepath string) error {
	p.fFilepath = filepath

	if !p.FSettings.IsValid() {
		return errors.NewError("load config settings")
	}

	if err := p.loadPubKeys(); err != nil {
		return errors.WrapError(err, "load public keys")
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
		if pubKey.GetSize() != pkg_settings.CAKeySize {
			return errors.NewError(fmt.Sprintf("not supported key size for '%s'", name))
		}
	}

	return nil
}

func (p *SConfig) GetNetwork() string {
	return p.FNetwork
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
	return p.FConnections
}

func (p *SConfig) GetService(name string) (string, bool) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	addr, ok := p.FServices[name]
	return addr, ok
}

func (p *STraffic) GetDownload() []string {
	if p == nil {
		return nil
	}
	return p.FDownload
}

func (p *STraffic) GetUpload() []string {
	if p == nil {
		return nil
	}
	return p.FUpload
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
	if p == nil {
		return ""
	}
	return p.FTCP
}

func (p *SAddress) GetHTTP() string {
	if p == nil {
		return ""
	}
	return p.FHTTP
}
