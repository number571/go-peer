package config

import (
	"os"
	"testing"
)

const (
	tcConfigFile   = "config_test.txt"
	tcLogging      = true
	tcWorkSizeBits = 20
	tcWaitTimeMS   = 1_000
	tcNetworkKey   = "network_key"
	tcAddress1     = "test_address1"
	tcAddress2     = "test_address2"
)

func testConfigDefaultInit(configPath string) {
	_, _ = BuildConfig(configPath, &SConfig{
		FSettings: &SConfigSettings{
			FWorkSizeBits: tcWorkSizeBits,
			FNetworkKey:   tcNetworkKey,
			FWaitTimeMS:   tcWaitTimeMS,
		},
		FLogging: []string{"info", "erro"},
		FConnection: &SConnection{
			FHLTHost: tcAddress1,
			FSrvHost: tcAddress2,
		},
	})
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

	if cfg.GetSettings().GetWaitTimeMS() != tcWaitTimeMS {
		t.Error("settings wait_time_ms is invalid")
		return
	}

	if cfg.GetSettings().GetNetworkKey() != tcNetworkKey {
		t.Error("settings network_key is invalid")
		return
	}

	if cfg.GetSettings().GetWorkSizeBits() != tcWorkSizeBits {
		t.Error("settings work_size_bits is invalid")
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

	if cfg.GetConnection().GetHLTHost() != tcAddress1 {
		t.Error("address hlt is invalid")
		return
	}

	if cfg.GetConnection().GetSrvHost() != tcAddress2 {
		t.Error("address srv is invalid")
		return
	}
}
