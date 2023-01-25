package config

import (
	"fmt"
	"sync"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/filesystem"
)

const (
	CLogInfo = "info"
	CLogWarn = "warn"
	CLogErro = "erro"
)

var (
	_ IConfig  = &SConfig{}
	_ IAddress = &SAddress{}
	_ ILogging = &sLogging{}
)

type SConfig struct {
	FNetwork string `json:"network,omitempty"`

	FAddress *SAddress `json:"address,omitempty"`
	FTraffic *STraffic `json:"traffic,omitempty"`
	FLogging []string  `json:"logging,omitempty"`

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

func NewConfig(filepath string, cfg *SConfig) (IConfig, error) {
	configFile := filesystem.OpenFile(filepath)

	if configFile.IsExist() {
		return nil, fmt.Errorf("config file '%s' already exist", filepath)
	}

	if err := configFile.Write(encoding.Serialize(cfg)); err != nil {
		return nil, err
	}

	if err := cfg.initConfig(filepath); err != nil {
		return nil, err
	}
	return cfg, nil
}

func LoadConfig(filepath string) (IConfig, error) {
	configFile := filesystem.OpenFile(filepath)

	if !configFile.IsExist() {
		return nil, fmt.Errorf("config file '%s' does not exist", filepath)
	}

	bytes, err := configFile.Read()
	if err != nil {
		return nil, err
	}

	cfg := new(SConfig)
	if err := encoding.Deserialize(bytes, cfg); err != nil {
		return nil, err
	}

	if err := cfg.initConfig(filepath); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (cfg *SConfig) initConfig(filepath string) error {
	cfg.fFilepath = filepath

	if err := cfg.loadPubKeys(); err != nil {
		return err
	}

	if err := cfg.loadLogging(); err != nil {
		return err
	}

	return nil
}

func (cfg *SConfig) loadLogging() error {
	// [info, warn, erro]
	logging := sLogging(make([]bool, 3))

	mapping := map[string]int{
		"info": 0,
		"warn": 1,
		"erro": 2,
	}

	for _, v := range cfg.FLogging {
		logType, ok := mapping[v]
		if !ok {
			return fmt.Errorf("undefined log type '%s'", v)
		}
		logging[logType] = true
	}

	cfg.fLogging = &logging
	return nil
}

func (cfg *SConfig) loadPubKeys() error {
	cfg.fFriends = make(map[string]asymmetric.IPubKey)

	mapping := make(map[string]interface{})
	for name, val := range cfg.FFriends {
		if _, ok := mapping[val]; ok {
			return fmt.Errorf("found public key duplicate '%s'", val)
		}
		pubKey := asymmetric.LoadRSAPubKey(val)
		if pubKey == nil {
			return fmt.Errorf("public key is nil for '%s'", name)
		}
		cfg.fFriends[name] = pubKey
		if pubKey.Size() != pkg_settings.CAKeySize {
			return fmt.Errorf("not supported key size for '%s'", name)
		}
	}

	return nil
}

func (cfg *SConfig) Network() string {
	return cfg.FNetwork
}

func (cfg *SConfig) Friends() map[string]asymmetric.IPubKey {
	cfg.fMutex.Lock()
	defer cfg.fMutex.Unlock()

	result := make(map[string]asymmetric.IPubKey)
	for k, v := range cfg.fFriends {
		result[k] = v
	}
	return result
}

func (cfg *SConfig) Logging() ILogging {
	return cfg.fLogging
}

func (cfg *SConfig) Address() IAddress {
	return cfg.FAddress
}

func (cfg *SConfig) Connections() []string {
	return cfg.FConnections
}

func (cfg *SConfig) Service(name string) (string, bool) {
	cfg.fMutex.Lock()
	defer cfg.fMutex.Unlock()

	addr, ok := cfg.FServices[name]
	return addr, ok
}

func (traffic *STraffic) Download() []string {
	if traffic == nil {
		return nil
	}
	return traffic.FDownload
}

func (traffic *STraffic) Upload() []string {
	if traffic == nil {
		return nil
	}
	return traffic.FUpload
}

func (logging *sLogging) Info() bool {
	return (*logging)[0]
}

func (logging *sLogging) Warn() bool {
	return (*logging)[1]
}

func (logging *sLogging) Erro() bool {
	return (*logging)[2]
}

func (address *SAddress) TCP() string {
	if address == nil {
		return ""
	}
	return address.FTCP
}

func (address *SAddress) HTTP() string {
	if address == nil {
		return ""
	}
	return address.FHTTP
}
