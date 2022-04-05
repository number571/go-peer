package config

import (
	"fmt"
	"sync"

	"github.com/number571/go-peer/cmd/hls/utils"
	"github.com/number571/go-peer/crypto"
)

var (
	_ IConfig = &sConfig{}
	_ iBlock  = &sBlock{}
)

type sBlock struct {
	FRedirect bool   `json:"redirect"`
	FAddress  string `json:"address"`
}

type sConfig struct {
	fMutex     sync.Mutex
	fPubKeys   []crypto.IPubKey
	FF2F       bool               `json:"f2f_mode"`
	FAddress   string             `json:"address"`
	FPubKeys   []string           `json:"pub_keys"`
	FConnects  []string           `json:"connects"`
	FServices  map[string]*sBlock `json:"services"`
	FCleanCron string             `json:"clean_cron"`
}

const (
	// friend-to-friend option
	cDefaultF2FMode = false

	// create local hls
	cDefaultAddress = "localhost:9571"

	// cron of clean database
	cDefaultCleanCron = "0 0 * * *"
)

var (
	// connect to another hls's
	gDefaultConnects = []string{
		"127.0.0.2:9571",
	}

	// another receivers of package
	gDefaultPubKeys = []string{
		`Pub(go-peer\rsa){30818902818100C709DA63096CEDBA0DD6B5DD9465B412268C00509757A8EBD9096E17BEEC17C25A3A8F246E1591554CD214F4B27254EFA811F8BE441A03B37B3C8B390484C74C2294A4C895AA925D723E0065A877D4502CC010996863821E7348348E4E96CDD4CB7A852B2E2853C8FDEE556C4F89F6C3295EAC00DAEE86DD94E25F9703F368C70203010001}`,
	}

	// crypto-address -> network-address
	gDefaultServices = map[string]*sBlock{
		"hidden-default-service": {
			FRedirect: true,
			FAddress:  "localhost:8080",
		},
	}
)

func NewConfig(filepath string) IConfig {
	var cfg = new(sConfig)

	if !utils.FileIsExist(filepath) {
		cfg = &sConfig{
			FF2F:       cDefaultF2FMode,
			FAddress:   cDefaultAddress,
			FConnects:  gDefaultConnects,
			FPubKeys:   gDefaultPubKeys,
			FServices:  gDefaultServices,
			FCleanCron: cDefaultCleanCron,
		}
		err := utils.WriteFile(filepath, utils.Serialize(cfg))
		if err != nil {
			panic(err)
		}
	} else {
		err := utils.Deserialize(utils.ReadFile(filepath), cfg)
		if err != nil {
			panic(err)
		}
	}

	for _, val := range cfg.FPubKeys {
		pubKey := crypto.LoadPubKey(val)
		if pubKey == nil {
			panic(fmt.Sprintf("public key is nil: '%s'", val))
		}
		cfg.fPubKeys = append(cfg.fPubKeys, pubKey)
	}

	return cfg
}

func (cfg *sConfig) F2F() bool {
	return cfg.FF2F
}

func (cfg *sConfig) Address() string {
	return cfg.FAddress
}

func (cfg *sConfig) PubKeys() []crypto.IPubKey {
	return cfg.fPubKeys
}

func (cfg *sConfig) Connections() []string {
	return cfg.FConnects
}

func (cfg *sConfig) GetService(name string) (iBlock, bool) {
	cfg.fMutex.Lock()
	defer cfg.fMutex.Unlock()

	addr, ok := cfg.FServices[name]
	return addr, ok
}

func (block *sBlock) Address() string {
	return block.FAddress
}

func (block *sBlock) IsRedirect() bool {
	return block.FRedirect
}

func (cfg *sConfig) CleanCron() string {
	return cfg.FCleanCron
}
