package config

import (
	"os"
	"testing"

	"github.com/number571/go-peer/utils"
)

const (
	tcConfigFile = "config_test.txt"
)

const (
	tcAddress    = "test_address"
	tcCleanCron  = "0 0 0 0 0"
	tcConnection = "test_connection"
)

const (
	tcConfig = `{
	"address": "test_address",
	"clean_cron": "0 0 0 0 0",
	"connections": [
		"test_connection"
	]
}`
)

func testConfigDefaultInit(configPath string) {
	utils.NewFile(configPath).Write([]byte(tcConfig))
}

func TestConfig(t *testing.T) {
	testConfigDefaultInit(tcConfigFile)
	defer os.Remove(tcConfigFile)

	cfg := NewConfig(tcConfigFile)

	if cfg.Address() != tcAddress {
		t.Error("address is invalid")
	}

	if cfg.CleanCron() != tcCleanCron {
		t.Error("clean_cron is invalid")
	}

	if cfg.Connections()[0] != tcConnection {
		t.Error("connections is invalid")
	}
}
