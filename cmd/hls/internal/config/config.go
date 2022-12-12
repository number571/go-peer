package config

import (
	"fmt"
	"sync"

	hls_settings "github.com/number571/go-peer/cmd/hls/internal/settings"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/encoding"
	"github.com/number571/go-peer/modules/filesystem"
)

var (
	_ IConfig  = &SConfig{}
	_ iAddress = &SAddress{}
)

type SConfig struct {
	FNetwork string `json:"network,omitempty"`

	FAddress  *SAddress         `json:"address,omitempty"`
	FServices map[string]string `json:"services,omitempty"`

	FConnections []string          `json:"connections,omitempty"`
	FFriends     map[string]string `json:"friends,omitempty"`

	fFilepath string
	fMutex    sync.Mutex
	fFriends  map[string]asymmetric.IPubKey
}

type SAddress struct {
	FTCP  string `json:"tcp,omitempty"`
	FHTTP string `json:"http,omitempty"`
}

type SKey struct {
	FStorage string `json:"storage,omitempty"`
	FNetwork string `json:"network,omitempty"`
}

func NewConfig(filepath string, cfg *SConfig) (IConfig, error) {
	configFile := filesystem.OpenFile(filepath)

	if configFile.IsExist() {
		return nil, fmt.Errorf("config file '%s' already exist", filepath)
	}

	if err := configFile.Write(encoding.Serialize(cfg)); err != nil {
		return nil, err
	}

	cfg.fFilepath = filepath
	if err := cfg.loadPubKeys(); err != nil {
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

	cfg.fFilepath = filepath
	if err := cfg.loadPubKeys(); err != nil {
		return nil, err
	}

	return cfg, nil
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
		if pubKey.Size() != hls_settings.CAKeySize {
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

func (cfg *SConfig) Address() iAddress {
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

func (key *SKey) Storage() string {
	if key == nil {
		return ""
	}
	return key.FStorage
}

func (key *SKey) Network() string {
	if key == nil {
		return ""
	}
	return key.FNetwork
}
