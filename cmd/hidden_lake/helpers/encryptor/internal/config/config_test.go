package config

import (
	"os"
	"testing"

	testutils "github.com/number571/go-peer/test/_data"
)

const (
	tcConfigFile  = "config_test.txt"
	tcLogging     = true
	tcNetwork     = "test_network"
	tcAddress1    = "test_address1"
	tcAddress2    = "test_address2"
	tcMessageSize = (1 << 20)
	tcWorkSize    = 22
)

func testConfigDefaultInit(configPath string) {
	_, _ = BuildConfig(configPath, &SConfig{
		FSettings: &SConfigSettings{
			FMessageSizeBytes: tcMessageSize,
			FWorkSizeBits:     tcWorkSize,
			FNetworkKey:       tcNetwork,
			FKeySizeBits:      testutils.TcKeySize,
		},
		FLogging: []string{"info", "erro"},
		FAddress: &SAddress{
			FHTTP:  tcAddress1,
			FPPROF: tcAddress2,
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

	if cfg.GetSettings().GetWorkSizeBits() != tcWorkSize {
		t.Error("settings work size is invalid")
		return
	}

	if cfg.GetSettings().GetMessageSizeBytes() != tcMessageSize {
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

	if cfg.GetSettings().GetNetworkKey() != tcNetwork {
		t.Error("network is invalid")
		return
	}

	if cfg.GetAddress().GetHTTP() != tcAddress1 {
		t.Error("address http is invalid")
		return
	}

	if cfg.GetAddress().GetPPROF() != tcAddress2 {
		t.Error("address pprof is invalid")
		return
	}
}
