package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/number571/go-peer/internal/language"
)

const (
	tcLogging    = true
	tcConfigFile = "config_test.txt"
)

const (
	tcConfigTemplate = `settings:
  messages_capacity: %d
  language: RUS
  pseudonym: '%s'
  storage_key: '%s'
logging:
  - info
  - erro
address:
  interface: '%s'
  incoming: '%s'
  pprof: '%s'
connection: '%s'`
)

const (
	tcPseudonym         = "Alice"
	tcAddressInterface  = "address_interface"
	tcAddressIncoming   = "address_incoming"
	tcAddressPPROF      = "address_pprof"
	tcConnectionService = "connection_service"
	tcStorageKey        = "storage_key"
	tcMessageSize       = (1 << 20)
	tcMessagesCapacity  = 1000
)

func testNewConfigString() string {
	return fmt.Sprintf(
		tcConfigTemplate,
		tcMessagesCapacity,
		tcPseudonym,
		tcStorageKey,
		tcAddressInterface,
		tcAddressIncoming,
		tcAddressPPROF,
		tcConnectionService,
	)
}

func testConfigDefaultInit(configPath string) {
	_ = os.WriteFile(configPath, []byte(testNewConfigString()), 0o644)
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

	if cfg.GetSettings().GetMessagesCapacity() != tcMessagesCapacity {
		t.Error("settings message capacity size is invalid")
		return
	}

	if cfg.GetSettings().GetLanguage() != language.CLangRUS {
		t.Error("settings language is invalid")
		return
	}

	if cfg.GetSettings().GetStorageKey() != tcStorageKey {
		t.Error("settings storage_key is invalid")
		return
	}

	if cfg.GetSettings().GetPseudonym() != tcPseudonym {
		t.Error("settings pseudonym is invalid")
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

}
