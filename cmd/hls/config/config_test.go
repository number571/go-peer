package config

import (
	"os"
	"testing"

	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/filesystem"
)

const (
	tcConfigFile = "config_test.txt"
)

const (
	tcAddressTCP  = "test_address_tcp"
	tcAddressHTTP = "test_address_http"
)

const (
	tcConfig = `{
	"address": {
		"tcp": "test_address_tcp",
		"http": "test_address_http"
	},
	"connections": [
		"test_connect1",
		"test_connect2",
		"test_connect3"
	],
	"friends": [
		"Pub(go-peer/rsa){30818902818100C709DA63096CEDBA0DD6B5DD9465B412268C00509757A8EBD9096E17BEEC17C25A3A8F246E1591554CD214F4B27254EFA811F8BE441A03B37B3C8B390484C74C2294A4C895AA925D723E0065A877D4502CC010996863821E7348348E4E96CDD4CB7A852B2E2853C8FDEE556C4F89F6C3295EAC00DAEE86DD94E25F9703F368C70203010001}"
	],
	"services": {
		"test_service1": "test_address1",
		"test_service2": "test_address2",
		"test_service3": "test_address3"
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
	tgServices = map[string]string{
		"test_service1": "test_address1",
		"test_service2": "test_address2",
		"test_service3": "test_address3",
	}
)

func testConfigDefaultInit(configPath string) {
	filesystem.OpenFile(configPath).Write([]byte(tcConfig))
}

func TestConfig(t *testing.T) {
	testConfigDefaultInit(tcConfigFile)
	defer os.Remove(tcConfigFile)

	cfg := NewConfig(tcConfigFile)

	if cfg.Address().TCP() != tcAddressTCP {
		t.Errorf("address_tcp is invalid")
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
		v1, ok := cfg.Service(k)
		if !ok {
			t.Errorf("service undefined '%s'", k)
		}
		if v != v1 {
			t.Errorf("service address is invalid '%s'", v1)
		}
	}

	for i, v := range tgPubKeys {
		v1 := cfg.Friends()[i]
		pubKey := asymmetric.LoadRSAPubKey(v)
		if pubKey.Address().String() != v1.Address().String() {
			t.Errorf("public key is invalid '%s'", v1)
		}
	}
}
