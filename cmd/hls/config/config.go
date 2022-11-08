package config

import (
	"fmt"
	"sync"

	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/encoding"
	"github.com/number571/go-peer/modules/filesystem"
)

var (
	_ IConfig  = &sConfig{}
	_ iAddress = &sAddress{}
)

type sConfig struct {
	FNetwork string `json:"network,omitempty"`

	FAddress  *sAddress         `json:"address,omitempty"`
	FServices map[string]string `json:"services,omitempty"`

	FConnections []string `json:"connections,omitempty"`
	FFriends     []string `json:"friends,omitempty"`

	fFilepath string
	fMutex    sync.Mutex
	fFriends  []asymmetric.IPubKey
}

type sAddress struct {
	FTCP  string `json:"tcp,omitempty"`
	FHTTP string `json:"http,omitempty"`
}

var (
	// network key
	cNetworkKey = "hls-network-key"

	// create local hls
	cDefaultAddress = &sAddress{
		"localhost:9571",
		"localhost:9572",
	}

	// connect to another hls's
	gDefaultConnects = []string{
		cDefaultAddress.FTCP,
	}

	// another receivers of package
	gDefaultPubKeys = []string{
		`Pub(go-peer/rsa){30818902818100C709DA63096CEDBA0DD6B5DD9465B412268C00509757A8EBD9096E17BEEC17C25A3A8F246E1591554CD214F4B27254EFA811F8BE441A03B37B3C8B390484C74C2294A4C895AA925D723E0065A877D4502CC010996863821E7348348E4E96CDD4CB7A852B2E2853C8FDEE556C4F89F6C3295EAC00DAEE86DD94E25F9703F368C70203010001}`,
	}

	// crypto-address -> network-address
	gDefaultServices = map[string]string{
		"hidden-default-service": "localhost:8080",
	}
)

func NewConfig(filepath string) IConfig {
	var cfg = new(sConfig)

	if !filesystem.OpenFile(filepath).IsExist() {
		cfg = &sConfig{
			FNetwork:     cNetworkKey,
			FAddress:     cDefaultAddress,
			FServices:    gDefaultServices,
			FFriends:     gDefaultPubKeys,
			FConnections: gDefaultConnects,
		}
		err := filesystem.OpenFile(filepath).Write(encoding.Serialize(cfg))
		if err != nil {
			panic(err)
		}
	} else {
		bytes, err := filesystem.OpenFile(filepath).Read()
		if err != nil {
			panic(err)
		}
		err = encoding.Deserialize(bytes, cfg)
		if err != nil {
			panic(err)
		}
	}

	cfg.fFilepath = filepath
	for _, val := range cfg.FFriends {
		pubKey := asymmetric.LoadRSAPubKey(val)
		if pubKey == nil {
			panic(fmt.Sprintf("public key is nil: '%s'", val))
		}
		cfg.fFriends = append(cfg.fFriends, pubKey)
	}

	return cfg
}

func (cfg *sConfig) Network() string {
	return cfg.FNetwork
}

func (cfg *sConfig) Friends() []asymmetric.IPubKey {
	return cfg.fFriends
}

func (cfg *sConfig) Address() iAddress {
	return cfg.FAddress
}

func (cfg *sConfig) Connections() []string {
	return cfg.FConnections
}

func (cfg *sConfig) Service(name string) (string, bool) {
	cfg.fMutex.Lock()
	defer cfg.fMutex.Unlock()

	addr, ok := cfg.FServices[name]
	return addr, ok
}

func (address *sAddress) TCP() string {
	if address == nil {
		return ""
	}
	return address.FTCP
}

func (address *sAddress) HTTP() string {
	if address == nil {
		return ""
	}
	return address.FHTTP
}
