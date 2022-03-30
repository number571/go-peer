package config

import (
	"testing"

	"github.com/number571/go-peer/crypto"
)

const (
	tcF2F        = true
	tcConfigFile = "config_test.txt"
	tcAddress    = "test_address"
	tcCleanCron  = "0 0 0 0 0"
)

var (
	tgConnects = []string{
		"test_connect1",
		"test_connect2",
		"test_connect3",
	}
	tgPubKeys = []string{
		`Pub(go-peer\rsa){30818902818100C709DA63096CEDBA0DD6B5DD9465B412268C00509757A8EBD9096E17BEEC17C25A3A8F246E1591554CD214F4B27254EFA811F8BE441A03B37B3C8B390484C74C2294A4C895AA925D723E0065A877D4502CC010996863821E7348348E4E96CDD4CB7A852B2E2853C8FDEE556C4F89F6C3295EAC00DAEE86DD94E25F9703F368C70203010001}`,
	}
	tgServices = map[string]*sBlock{
		"test_service1": &sBlock{false, "test_address1"},
		"test_service2": &sBlock{true, "test_address2"},
		"test_service3": &sBlock{true, "test_address3"},
	}
)

func TestConfig(t *testing.T) {
	cfg := NewConfig(tcConfigFile)

	if cfg.F2F() != tcF2F {
		t.Errorf("f2f_mode is invalid")
	}

	if cfg.Address() != tcAddress {
		t.Errorf("address is invalid")
	}

	for i, v := range cfg.Connections() {
		if v != tgConnects[i] {
			t.Errorf("connection '%d' is invalid", i)
		}
	}

	for k, v := range tgServices {
		v1, ok := cfg.GetService(k)
		if !ok {
			t.Errorf("service undefined '%s'", k)
		}
		if v.IsRedirect() != v1.IsRedirect() {
			t.Errorf("service redirect is invalid '%s'", v1)
		}
		if v.Address() != v1.Address() {
			t.Errorf("service address is invalid '%s'", v1)
		}
	}

	for i, v := range tgPubKeys {
		v1 := cfg.PubKeys()[i]
		pubKey := crypto.LoadPubKey(v)
		if pubKey.Address() != v1.Address() {
			t.Errorf("public key is invalid '%s'", v1)
		}
	}

	if cfg.CleanCron() != tcCleanCron {
		t.Errorf("clean_cron is invalid")
	}
}
