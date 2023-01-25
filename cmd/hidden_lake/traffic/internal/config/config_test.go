package config

import (
	"os"
	"testing"
)

const (
	tcConfigFile = "config_test.txt"
	tcAddress    = "test_address"
	tcConnection = "test_connection"
)

func testConfigDefaultInit(configPath string) {
	_, _ = NewConfig(configPath, &SConfig{
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

	if cfg.Address() != tcAddress {
		t.Error("address is invalid")
	}

	if cfg.Connection() != tcConnection {
		t.Error("connection is invalid")
	}
}
