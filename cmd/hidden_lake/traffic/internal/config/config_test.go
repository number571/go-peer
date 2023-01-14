package config

import (
	"os"
	"testing"
)

const (
	tcConfigFile = "config_test.txt"
	tcAddress    = "test_address"
)

func testConfigDefaultInit(configPath string) {
	_, _ = NewConfig(configPath, &SConfig{
		FAddress: tcAddress,
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
}
