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
  page_offset: %d
  retry_num: %d
  language: RUS
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
	tcAddressInterface  = "address_interface"
	tcAddressIncoming   = "address_incoming"
	tcAddressPPROF      = "address_pprof"
	tcConnectionService = "connection_service"
	tcMessageSize       = (1 << 20)
	tcPageOffset        = 10
	tcRetryNum          = 2
)

func testNewConfigString() string {
	return fmt.Sprintf(
		tcConfigTemplate,
		tcPageOffset,
		tcRetryNum,
		tcAddressInterface,
		tcAddressIncoming,
		tcAddressPPROF,
		tcConnectionService,
	)
}

func testConfigDefaultInit(configPath string) {
	os.WriteFile(configPath, []byte(testNewConfigString()), 0o644)
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

	if cfg.GetSettings().GetPageOffset() != tcPageOffset {
		t.Error("settings.page_offset is invalid")
		return
	}

	if cfg.GetSettings().GetRetryNum() != tcRetryNum {
		t.Error("settings.retry_num is invalid")
		return
	}

	if cfg.GetSettings().GetLanguage() != language.CLangRUS {
		t.Error("settings language is invalid")
		return
	}
}
