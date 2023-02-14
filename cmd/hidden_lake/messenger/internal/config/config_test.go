package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/filesystem"
)

const (
	tcConfigFile = "config_test.txt"
)

const (
	tcConfigTemplate = `{
	"address": {
		"interface": "%s",
		"incoming": "%s"
	},
	"connection": {
		"service": "%s",
		"traffic": "%s"
	},
	"storage_key": "%s"
}`
)

const (
	tcAddressInterface  = "address_interface"
	tcAddressIncoming   = "address_incoming"
	tcConnectionService = "connection_service"
	tcConnectionTraffic = "connection_traffic"
	tcStorageKey        = "storage_key"
)

func testNewConfigString() string {
	return fmt.Sprintf(
		tcConfigTemplate,
		tcAddressInterface,
		tcAddressIncoming,
		tcConnectionService,
		tcConnectionTraffic,
		tcStorageKey,
	)
}

func testConfigDefaultInit(configPath string) {
	filesystem.OpenFile(configPath).Write([]byte(testNewConfigString()))
}

func TestConfig(t *testing.T) {
	testConfigDefaultInit(tcConfigFile)
	defer os.Remove(tcConfigFile)

	cfg, err := LoadConfig(tcConfigFile)
	if err != nil {
		t.Error(err)
		return
	}

	if cfg.Address().Interface() != tcAddressInterface {
		t.Error("address.interface is invalid")
	}

	if cfg.Address().Incoming() != tcAddressIncoming {
		t.Error("address.incoming is invalid")
	}

	if cfg.Connection().Service() != tcConnectionService {
		t.Error("connection.service is invalid")
	}

	if cfg.Connection().Traffic() != tcConnectionTraffic {
		t.Error("connection.traffic is invalid")
	}

	if cfg.StorageKey() != tcStorageKey {
		t.Error("storage_key is invalid")
	}
}
