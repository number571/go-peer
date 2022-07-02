package config

import (
	"os"
	"testing"

	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/utils"
)

const (
	tcConfigFile = "config_test.txt"
)

const (
	tcOnlineChecker = true
	tcF2F           = true
	tcAddressHLS    = "test_address_hls"
	tcAddressHTTP   = "test_address_http"
	tcCleanCron     = "0 0 0 0 0"
)

const (
	tcConfig = `{
	"clean_cron": "0 0 0 0 0",
	"address": {
		"hls": "test_address_hls",
		"http": "test_address_http"
	},
	"connections": [
		"test_connect1",
		"test_connect2",
		"test_connect3"
	],
	"f2f_mode": {
		"status": true,
		"pub_keys": [
			"Pub(go-peer/rsa){30818902818100C709DA63096CEDBA0DD6B5DD9465B412268C00509757A8EBD9096E17BEEC17C25A3A8F246E1591554CD214F4B27254EFA811F8BE441A03B37B3C8B390484C74C2294A4C895AA925D723E0065A877D4502CC010996863821E7348348E4E96CDD4CB7A852B2E2853C8FDEE556C4F89F6C3295EAC00DAEE86DD94E25F9703F368C70203010001}"
		]
	},
	"online_checker": {
		"status": true,
		"pub_keys": [
			"Pub(go-peer/rsa){30818902818100C709DA63096CEDBA0DD6B5DD9465B412268C00509757A8EBD9096E17BEEC17C25A3A8F246E1591554CD214F4B27254EFA811F8BE441A03B37B3C8B390484C74C2294A4C895AA925D723E0065A877D4502CC010996863821E7348348E4E96CDD4CB7A852B2E2853C8FDEE556C4F89F6C3295EAC00DAEE86DD94E25F9703F368C70203010001}"
		]
	},
	"services": {
		"test_service1": {
			"redirect": false,
			"address": "test_address1"
		},
		"test_service2": {
			"redirect": true,
			"address": "test_address2"
		},
		"test_service3": {
			"redirect": true,
			"address": "test_address3"
		}
	}
}`
)

var (
	tgConnects = []string{
		"test_connect1",
		"test_connect2",
		"test_connect3",
	}
	tgPubKeys = []string{
		`Pub(go-peer/rsa){30818902818100C709DA63096CEDBA0DD6B5DD9465B412268C00509757A8EBD9096E17BEEC17C25A3A8F246E1591554CD214F4B27254EFA811F8BE441A03B37B3C8B390484C74C2294A4C895AA925D723E0065A877D4502CC010996863821E7348348E4E96CDD4CB7A852B2E2853C8FDEE556C4F89F6C3295EAC00DAEE86DD94E25F9703F368C70203010001}`,
	}
	tgServices = map[string]*sBlock{
		"test_service1": {false, "test_address1"},
		"test_service2": {true, "test_address2"},
		"test_service3": {true, "test_address3"},
	}
)

func testConfigDefaultInit(configPath string) {
	utils.OpenFile(configPath).Write([]byte(tcConfig))
}

func TestConfig(t *testing.T) {
	testConfigDefaultInit(tcConfigFile)
	defer os.Remove(tcConfigFile)

	cfg := NewConfig(tcConfigFile)

	if cfg.F2F().Status() != tcF2F {
		t.Errorf("f2f_mode is invalid")
	}

	if cfg.OnlineChecker().Status() != tcOnlineChecker {
		t.Errorf("online_checker is invalid")
	}

	if cfg.Address().HLS() != tcAddressHLS {
		t.Errorf("address_hls is invalid")
	}

	if cfg.Address().HTTP() != tcAddressHTTP {
		t.Errorf("address_http is invalid")
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
		v1 := cfg.F2F().PubKeys()[i]
		pubKey := asymmetric.LoadRSAPubKey(v)
		if pubKey.Address().String() != v1.Address().String() {
			t.Errorf("public key is invalid '%s'", v1)
		}

		v2 := cfg.OnlineChecker().PubKeys()[i]
		if pubKey.Address().String() != v2.Address().String() {
			t.Errorf("public key is invalid '%s'", v2)
		}
	}

	if cfg.CleanCron() != tcCleanCron {
		t.Errorf("clean_cron is invalid")
	}
}
