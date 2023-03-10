package config

import (
	"os"
	"testing"
)

const (
	tcConfigFile = "config_test.txt"
	tcLogging    = true
	tcNetwork    = "test_network"
	tcAddress    = "test_address"
	tcConnection = "test_connection"
)

func testConfigDefaultInit(configPath string) {
	_, _ = BuildConfig(configPath, &SConfig{
		FLogging:    []string{"info", "erro"},
		FNetwork:    tcNetwork,
		FAddress:    tcAddress,
		FConnection: tcConnection,
	})
}

func TestConfig(t *testing.T) {
	testConfigDefaultInit(tcConfigFile)
	defer os.Remove(tcConfigFile)

	cfg, err := LoadConfig(tcConfigFile)
	if err != nil {
		t.Error(err)
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

	if cfg.GetNetwork() != tcNetwork {
		t.Error("network is invalid")
	}

	if cfg.GetAddress() != tcAddress {
		t.Error("address is invalid")
	}

	if cfg.GetConnection() != tcConnection {
		t.Error("connection is invalid")
	}
}
