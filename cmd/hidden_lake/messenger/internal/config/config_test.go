package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/filesystem"
)

const (
	tcLogging    = true
	tcConfigFile = "config_test.txt"
)

const (
	tcConfigTemplate = `{
	"settings": {
		"message_size": %d,
		"work_size": %d
	},
	"logging": ["info", "erro"],
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
	tcMessageSize       = (1 << 20)
	tcWorkSize          = 20
)

func testNewConfigString() string {
	return fmt.Sprintf(
		tcConfigTemplate,
		tcMessageSize,
		tcWorkSize,
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

	if cfg.GetWorkSize() != tcWorkSize {
		t.Error("settings work size is invalid")
		return
	}

	if cfg.GetMessageSize() != tcMessageSize {
		t.Error("settings message size is invalid")
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

	if cfg.GetAddress().GetInterface() != tcAddressInterface {
		t.Error("address.interface is invalid")
	}

	if cfg.GetAddress().GetIncoming() != tcAddressIncoming {
		t.Error("address.incoming is invalid")
	}

	if cfg.GetConnection().GetService() != tcConnectionService {
		t.Error("connection.service is invalid")
	}

	if cfg.GetConnection().GetTraffic() != tcConnectionTraffic {
		t.Error("connection.traffic is invalid")
	}

	if cfg.GetStorageKey() != tcStorageKey {
		t.Error("storage_key is invalid")
	}
}
