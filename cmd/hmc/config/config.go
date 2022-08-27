package config

import (
	"fmt"

	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/encoding"
	"github.com/number571/go-peer/modules/filesystem"
)

var (
	_ IConfig = &sConfig{}
)

type sConfig struct {
	FF2F         bool       `json:"f2f_mode"`
	FConnections []string   `json:"connections"`
	FFriends     []*sFriend `json:"friends"`
	fFriends     []iFriend
}

type sFriend struct {
	FName   string `json:"name"`
	FPubKey string `json:"pub_key"`
	fPubKey asymmetric.IPubKey
}

const (
	// create local hms
	cDefaultConnection = "http://localhost:8080"

	// example of friend's name
	gDefaultName = "default"

	// example of friend's public key
	gDefaultPubKey = `Pub(go-peer/rsa){30818902818100C709DA63096CEDBA0DD6B5DD9465B412268C00509757A8EBD9096E17BEEC17C25A3A8F246E1591554CD214F4B27254EFA811F8BE441A03B37B3C8B390484C74C2294A4C895AA925D723E0065A877D4502CC010996863821E7348348E4E96CDD4CB7A852B2E2853C8FDEE556C4F89F6C3295EAC00DAEE86DD94E25F9703F368C70203010001}`
)

func NewConfig(filepath string) IConfig {
	var cfg = new(sConfig)

	if !filesystem.OpenFile(filepath).IsExist() {
		cfg = &sConfig{
			FF2F:         false,
			FConnections: []string{cDefaultConnection},
			FFriends: []*sFriend{{
				FName:   gDefaultName,
				FPubKey: gDefaultPubKey,
			}},
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

	for _, val := range cfg.FFriends {
		val.fPubKey = asymmetric.LoadRSAPubKey(val.FPubKey)
		if val.fPubKey == nil {
			panic(fmt.Sprintf("public key is nil: '%s'", val.FName))
		}
		cfg.fFriends = append(cfg.fFriends, val)
	}

	return cfg
}

func (cfg *sConfig) F2F() bool {
	return cfg.FF2F
}

func (cfg *sConfig) Connections() []string {
	return cfg.FConnections
}

func (cfg *sConfig) Friends() []iFriend {
	return cfg.fFriends
}

func (f *sFriend) Name() string {
	return f.FName
}

func (f *sFriend) PubKey() asymmetric.IPubKey {
	return f.fPubKey
}
