package config

import (
	"fmt"
	"sync"

	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/encoding"
	"github.com/number571/go-peer/utils"
)

var (
	_ IConfig  = &sConfig{}
	_ iAddress = &sAddress{}
	_ iBlock   = &sBlock{}
)

type sConfig struct {
	FAddress     *sAddress          `json:"address"`
	FConnections []string           `json:"connections"`
	FFriends     []string           `json:"friends"`
	FServices    map[string]*sBlock `json:"services"`

	fMutex   sync.Mutex
	fFriends []asymmetric.IPubKey
}

type sAddress struct {
	FTCP  string `json:"tcp"`
	FHTTP string `json:"http"`
}

type sBlock struct {
	FRedirect bool   `json:"redirect"`
	FAddress  string `json:"address"`
}

var (
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
	gDefaultServices = map[string]*sBlock{
		"hidden-default-service": {
			FRedirect: false,
			FAddress:  "localhost:8080",
		},
	}
)

func NewConfig(filepath string) IConfig {
	var cfg = new(sConfig)

	if !utils.OpenFile(filepath).IsExist() {
		cfg = &sConfig{
			FAddress:     cDefaultAddress,
			FFriends:     gDefaultPubKeys,
			FConnections: gDefaultConnects,
			FServices:    gDefaultServices,
		}
		err := utils.OpenFile(filepath).Write(encoding.Serialize(cfg))
		if err != nil {
			panic(err)
		}
	} else {
		bytes, err := utils.OpenFile(filepath).Read()
		if err != nil {
			panic(err)
		}
		err = encoding.Deserialize(bytes, cfg)
		if err != nil {
			panic(err)
		}
	}

	for _, val := range cfg.FFriends {
		pubKey := asymmetric.LoadRSAPubKey(val)
		if pubKey == nil {
			panic(fmt.Sprintf("public key is nil: '%s'", val))
		}
		cfg.fFriends = append(cfg.fFriends, pubKey)
	}

	return cfg
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

func (cfg *sConfig) GetService(name string) (iBlock, bool) {
	cfg.fMutex.Lock()
	defer cfg.fMutex.Unlock()

	addr, ok := cfg.FServices[name]
	return addr, ok
}

func (address *sAddress) TCP() string {
	return address.FTCP
}

func (address *sAddress) HTTP() string {
	return address.FHTTP
}

func (block *sBlock) Address() string {
	return block.FAddress
}

func (block *sBlock) IsRedirect() bool {
	return block.FRedirect
}
