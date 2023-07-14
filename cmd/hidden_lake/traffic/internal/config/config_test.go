package config

import (
	"os"
	"testing"

	"github.com/number571/go-peer/internal/settings"
)

const (
	tcConfigFile  = "config_test.txt"
	tcLogging     = true
	tcStorage     = true
	tcNetwork     = "test_network"
	tcAddress1    = "test_address1"
	tcAddress2    = "test_address2"
	tcConnection1 = "test_connection1"
	tcConnection2 = "test_connection2"
	tcConsumer1   = "test_consumer1"
	tcConsumer2   = "test_consumer2"
	tcMessageSize = (1 << 20)
	tcWorkSize    = 20
	tcCapMessages = 1000
)

func testConfigDefaultInit(configPath string) {
	_, _ = BuildConfig(configPath, &SConfig{
		SConfigSettings: settings.SConfigSettings{
			FSettings: settings.SConfigSettingsBlock{
				FMessageSizeBytes: tcMessageSize,
				FWorkSizeBits:     tcWorkSize,
				FMessagesCapacity: tcCapMessages,
			},
		},
		FLogging: []string{"info", "erro"},
		FNetwork: tcNetwork,
		FAddress: &SAddress{
			FTCP:  tcAddress1,
			FHTTP: tcAddress2,
		},
		FConnections: []string{
			tcConnection1,
			tcConnection2,
		},
		FConsumers: []string{
			tcConsumer1,
			tcConsumer2,
		},
	})
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

	if cfg.GetMessagesCapacity() != tcCapMessages {
		t.Error("settings cap messages is invalid")
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

	if cfg.GetNetwork() != tcNetwork {
		t.Error("network is invalid")
		return
	}

	if cfg.GetAddress().GetTCP() != tcAddress1 {
		t.Error("address tcp is invalid")
		return
	}

	if cfg.GetAddress().GetHTTP() != tcAddress2 {
		t.Error("address http is invalid")
		return
	}

	if len(cfg.GetConnections()) != 2 {
		t.Error("length of connections != 2")
		return
	}

	if cfg.GetConnections()[0] != tcConnection1 {
		t.Error("connection is invalid")
		return
	}

	if len(cfg.GetConsumers()) != 2 {
		t.Error("length of consumers != 2")
		return
	}

	if cfg.GetConsumers()[0] != tcConsumer1 {
		t.Error("consumers is invalid")
		return
	}
}
