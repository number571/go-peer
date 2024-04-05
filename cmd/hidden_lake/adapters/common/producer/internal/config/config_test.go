package config

import (
	"os"
	"testing"
)

const (
	tcConfigFile   = "config_test.txt"
	tcLogging      = true
	tcWorkSizeBits = 20
	tcNetworkKey   = "network_key"
	tcAddress1     = "test_address1"
	tcAddress2     = "test_address2"
)

func testConfigDefaultInit(configPath string) {
	_, _ = BuildConfig(configPath, &SConfig{
		FSettings: &SConfigSettings{
			FWorkSizeBits: tcWorkSizeBits,
			FNetworkKey:   tcNetworkKey,
		},
		FLogging: []string{"info", "erro"},
		FConnection: &SConnection{
			FSrvHost: tcAddress1,
		},
		FAddress: tcAddress2,
	})
}

func TestError(t *testing.T) {
	str := "value"
	err := &SConfigError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestConfig(t *testing.T) {
	t.Parallel()

	testConfigDefaultInit(tcConfigFile)
	defer os.Remove(tcConfigFile)

	cfg, err := LoadConfig(tcConfigFile)
	if err != nil {
		t.Error(err)
		return
	}

	if cfg.GetSettings().GetNetworkKey() != tcNetworkKey {
		t.Error("settings network_key is invalid")
		return
	}

	if cfg.GetSettings().GetWorkSizeBits() != tcWorkSizeBits {
		t.Error("settings work_size_bits is invalid")
		return
	}

	if cfg.GetLogging().HasInfo() != tcLogging {
		t.Error("logging.info is invalid")
		return
	}

	if cfg.GetLogging().HasErro() != tcLogging {
		t.Error("logging.erro is invalid")
		return
	}

	if cfg.GetLogging().HasWarn() == tcLogging {
		t.Error("logging.warn is invalid")
		return
	}

	if cfg.GetConnection().GetSrvHost() != tcAddress1 {
		t.Error("address connection is invalid")
		return
	}

	if cfg.GetAddress() != tcAddress2 {
		t.Error("address address is invalid")
		return
	}
}
