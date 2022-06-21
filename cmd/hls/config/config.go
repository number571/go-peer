package config

import (
	"fmt"
	"sync"

	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/encoding"
	"github.com/number571/go-peer/utils"
)

var (
	_ IConfig        = &sConfig{}
	_ iF2F           = &sF2F{}
	_ iOnlineChecker = &sOnlineChecker{}
	_ iAddress       = &sAddress{}
	_ iBlock         = &sBlock{}
)

type sConfig struct {
	fMutex         sync.Mutex
	FCleanCron     string             `json:"clean_cron"`
	FAddress       *sAddress          `json:"address"`
	FConnections   []string           `json:"connections"`
	FF2F           *sF2F              `json:"f2f_mode"`
	FOnlineChecker *sOnlineChecker    `json:"online_checker"`
	FServices      map[string]*sBlock `json:"services"`
}

type sF2F struct {
	FStatus  bool `json:"status"`
	fPubKeys []asymmetric.IPubKey
	FPubKeys []string `json:"pub_keys"`
}

type sOnlineChecker struct {
	FStatus  bool `json:"status"`
	fPubKeys []asymmetric.IPubKey
	FPubKeys []string `json:"pub_keys"`
}

type sAddress struct {
	FHLS  string `json:"hls"`
	FHTTP string `json:"http"`
}

type sBlock struct {
	FRedirect bool   `json:"redirect"`
	FAddress  string `json:"address"`
}

const (
	// cron of clean database
	cDefaultCleanCron = "0 0 * * *"
)

var (
	// friend-to-friend option
	cDefaultF2FMode = &sF2F{
		FStatus:  true,
		FPubKeys: gDefaultPubKeys,
	}

	// online checker option
	cDefaultOnlineChecker = &sOnlineChecker{
		FStatus:  true,
		FPubKeys: gDefaultPubKeys,
	}

	// create local hls
	cDefaultAddress = &sAddress{
		"localhost:9571",
		"localhost:9572",
	}

	// connect to another hls's
	gDefaultConnects = []string{
		cDefaultAddress.FHLS,
	}

	// another receivers of package
	gDefaultPubKeys = []string{
		`Pub(go-peer/rsa){30818902818100C709DA63096CEDBA0DD6B5DD9465B412268C00509757A8EBD9096E17BEEC17C25A3A8F246E1591554CD214F4B27254EFA811F8BE441A03B37B3C8B390484C74C2294A4C895AA925D723E0065A877D4502CC010996863821E7348348E4E96CDD4CB7A852B2E2853C8FDEE556C4F89F6C3295EAC00DAEE86DD94E25F9703F368C70203010001}`,
	}

	// crypto-address -> network-address
	gDefaultServices = map[string]*sBlock{
		"hidden-default-service": {
			FRedirect: true,
			FAddress:  "localhost:8571",
		},
	}
)

func NewConfig(filepath string) IConfig {
	var cfg = new(sConfig)

	if !utils.NewFile(filepath).IsExist() {
		cfg = &sConfig{
			FCleanCron:     cDefaultCleanCron,
			FAddress:       cDefaultAddress,
			FF2F:           cDefaultF2FMode,
			FConnections:   gDefaultConnects,
			FOnlineChecker: cDefaultOnlineChecker,
			FServices:      gDefaultServices,
		}
		err := utils.NewFile(filepath).Write(encoding.Serialize(cfg))
		if err != nil {
			panic(err)
		}
	} else {
		bytes, err := utils.NewFile(filepath).Read()
		if err != nil {
			panic(err)
		}
		err = encoding.Deserialize(bytes, cfg)
		if err != nil {
			panic(err)
		}
	}

	for _, val := range cfg.FOnlineChecker.FPubKeys {
		pubKey := asymmetric.LoadRSAPubKey(val)
		if pubKey == nil {
			panic(fmt.Sprintf("public key is nil: '%s'", val))
		}
		cfg.FOnlineChecker.fPubKeys = append(cfg.FOnlineChecker.fPubKeys, pubKey)
	}

	for _, val := range cfg.FF2F.FPubKeys {
		pubKey := asymmetric.LoadRSAPubKey(val)
		if pubKey == nil {
			panic(fmt.Sprintf("public key is nil: '%s'", val))
		}
		cfg.FF2F.fPubKeys = append(cfg.FF2F.fPubKeys, pubKey)
	}

	return cfg
}

func (cfg *sConfig) F2F() iF2F {
	return cfg.FF2F
}

func (f2f *sF2F) Status() bool {
	return f2f.FStatus
}

func (f2f *sF2F) PubKeys() []asymmetric.IPubKey {
	return f2f.fPubKeys
}

func (cfg *sConfig) Address() iAddress {
	return cfg.FAddress
}

func (cfg *sConfig) OnlineChecker() iOnlineChecker {
	return cfg.FOnlineChecker
}

func (onl *sOnlineChecker) Status() bool {
	return onl.FStatus
}

func (onl *sOnlineChecker) PubKeys() []asymmetric.IPubKey {
	return onl.fPubKeys
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

func (cfg *sConfig) CleanCron() string {
	return cfg.FCleanCron
}

func (address *sAddress) HLS() string {
	return address.FHLS
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
