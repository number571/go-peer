package config

import (
	"os"
	"testing"

	"github.com/number571/go-peer/cmd/hms/utils"
)

const (
	tcConfigFile = "config_test.txt"
)

const (
	tcAddress   = "test_address"
	tcCleanCron = "0 0 0 0 0"
)

const (
	tcConfig = `{
	"address": "test_address",
	"clean_cron": "0 0 0 0 0"
}`
)

func testConfigDefaultInit(configPath string) {
	utils.WriteFile(configPath, []byte(tcConfig))
}

func TestConfig(t *testing.T) {
	testConfigDefaultInit(tcConfigFile)
	defer os.Remove(tcConfigFile)

	cfg := NewConfig(tcConfigFile)

	if cfg.Address() != tcAddress {
		t.Errorf("address is invalid")
	}

	if cfg.CleanCron() != tcCleanCron {
		t.Errorf("clean_cron is invalid")
	}
}
