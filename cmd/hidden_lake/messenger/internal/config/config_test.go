package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/utils"
	"github.com/number571/go-peer/pkg/file_system"
)

const (
	tcLogging    = true
	tcConfigFile = "config_test.txt"
)

const (
	tcConfigTemplate = `{
	"settings": {
		"messages_capacity": %d
	},
	"logging": ["info", "erro"],
	"language": "RUS",
	"address": {
		"interface": "%s",
		"incoming": "%s",
		"pprof": "%s"
	},
	"connection": "%s",
	"storage_key": "%s"
}`
)

const (
	tcAddressInterface  = "address_interface"
	tcAddressIncoming   = "address_incoming"
	tcAddressPPROF      = "address_pprof"
	tcConnectionService = "connection_service"
	tcStorageKey        = "storage_key"
	tcMessageSize       = (1 << 20)
	tcKeySize           = 1024
	tcMessagesCapacity  = 1000
)

func testNewConfigString() string {
	return fmt.Sprintf(
		tcConfigTemplate,
		tcMessagesCapacity,
		tcAddressInterface,
		tcAddressIncoming,
		tcAddressPPROF,
		tcConnectionService,
		tcStorageKey,
	)
}

func testConfigDefaultInit(configPath string) {
	file_system.OpenFile(configPath).Write([]byte(testNewConfigString()))
}

func TestConfig(t *testing.T) {
	testConfigDefaultInit(tcConfigFile)
	defer os.Remove(tcConfigFile)

	cfg, err := LoadConfig(tcConfigFile)
	if err != nil {
		t.Error(err)
		return
	}

	if cfg.GetSettings().GetMessagesCapacity() != tcMessagesCapacity {
		t.Error("settings key size is invalid")
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

	if cfg.GetLanguage() != utils.CLangRUS {
		t.Error("language is invalid")
		return
	}

	if cfg.GetAddress().GetInterface() != tcAddressInterface {
		t.Error("address.interface is invalid")
		return
	}

	if cfg.GetAddress().GetIncoming() != tcAddressIncoming {
		t.Error("address.incoming is invalid")
		return
	}

	if cfg.GetAddress().GetPPROF() != tcAddressPPROF {
		t.Error("address.pprof is invalid")
		return
	}

	if cfg.GetConnection() != tcConnectionService {
		t.Error("connection.service is invalid")
		return
	}

	if cfg.GetStorageKey() != tcStorageKey {
		t.Error("storage_key is invalid")
		return
	}
}
