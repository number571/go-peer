package config

import (
	"os"
	"testing"

	"github.com/number571/go-peer/modules/filesystem"
)

const (
	tcConfigFile = "config_test.txt"
)

const (
	tcF2F        = true
	tcConnection = "test_connection"
	tcName       = "default"
	tcPubKey     = "Pub(go-peer/rsa){30818902818100C709DA63096CEDBA0DD6B5DD9465B412268C00509757A8EBD9096E17BEEC17C25A3A8F246E1591554CD214F4B27254EFA811F8BE441A03B37B3C8B390484C74C2294A4C895AA925D723E0065A877D4502CC010996863821E7348348E4E96CDD4CB7A852B2E2853C8FDEE556C4F89F6C3295EAC00DAEE86DD94E25F9703F368C70203010001}"
)

const (
	tcConfig = `{
	"f2f_mode": true,
	"connections": [
		"test_connection"
	],
	"friends": [
		{
			"name": "default",
			"pub_key": "Pub(go-peer/rsa){30818902818100C709DA63096CEDBA0DD6B5DD9465B412268C00509757A8EBD9096E17BEEC17C25A3A8F246E1591554CD214F4B27254EFA811F8BE441A03B37B3C8B390484C74C2294A4C895AA925D723E0065A877D4502CC010996863821E7348348E4E96CDD4CB7A852B2E2853C8FDEE556C4F89F6C3295EAC00DAEE86DD94E25F9703F368C70203010001}"
		}
	]
}`
)

func testConfigDefaultInit(configPath string) {
	filesystem.OpenFile(configPath).Write([]byte(tcConfig))
}

func TestConfig(t *testing.T) {
	testConfigDefaultInit(tcConfigFile)
	defer os.Remove(tcConfigFile)

	cfg := NewConfig(tcConfigFile)

	if cfg.F2F() != tcF2F {
		t.Error("f2f is invalid")
	}

	if cfg.Connections()[0] != tcConnection {
		t.Error("connections is invalid")
	}

	if cfg.Friends()[0].Name() != tcName {
		t.Error("name is invalid")
	}

	if cfg.Friends()[0].PubKey().String() != tcPubKey {
		t.Error("public key is invalid")
	}
}
