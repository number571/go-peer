package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/utils"
	"github.com/number571/go-peer/pkg/filesystem"
)

const (
	tcLogging    = true
	tcConfigFile = "config_test.txt"
)

const (
	tcConfigTemplate = `{
	"settings": {
		"message_size_bytes": %d,
		"work_size_bits": %d,
		"key_size_bits": %d,
		"messages_capacity": %d
	},
	"logging": ["info", "erro"],
	"language": "RUS",
	"address": {
		"interface": "%s",
		"incoming": "%s"
	},
	"connection": "%s",
	"backup_connections": [
		"%s"
	],
	"storage_key": "%s"
}`
)

const (
	tcAddressInterface  = "address_interface"
	tcAddressIncoming   = "address_incoming"
	tcConnectionService = "connection_service"
	tcConnectionBackup  = "connection_backup"
	tcStorageKey        = "storage_key"
	tcMessageSize       = (1 << 20)
	tcWorkSize          = 20
	tcKeySize           = 1024
	tcMessagesCapacity  = 1000
)

func testNewConfigString() string {
	return fmt.Sprintf(
		tcConfigTemplate,
		tcMessageSize,
		tcWorkSize,
		tcKeySize,
		tcMessagesCapacity,
		tcAddressInterface,
		tcAddressIncoming,
		tcConnectionService,
		tcConnectionBackup,
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

	if cfg.GetWorkSizeBits() != tcWorkSize {
		t.Error("settings work size is invalid")
		return
	}

	if cfg.GetMessageSizeBytes() != tcMessageSize {
		t.Error("settings message size is invalid")
		return
	}

	if cfg.GetKeySizeBits() != tcKeySize {
		t.Error("settings key size is invalid")
		return
	}

	if cfg.GetMessagesCapacity() != tcMessagesCapacity {
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

	if cfg.GetConnection() != tcConnectionService {
		t.Error("connection.service is invalid")
		return
	}

	if len(cfg.GetBackupConnections()) != 1 {
		t.Error("length of connections.backup is invalid")
		return
	}

	if cfg.GetBackupConnections()[0] != tcConnectionBackup {
		t.Error("connections[0].backup is invalid")
		return
	}

	if cfg.GetStorageKey() != tcStorageKey {
		t.Error("storage_key is invalid")
		return
	}
}
