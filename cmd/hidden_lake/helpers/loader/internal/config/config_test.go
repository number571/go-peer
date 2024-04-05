package config

import (
	"os"
	"testing"
)

const (
	tcConfigFile  = "config_test.txt"
	tcLogging     = true
	tcNetwork     = "test_network"
	tcAddress1    = "test_address1"
	tcAddress2    = "test_address2"
	tcProducer1   = "test_producer1"
	tcProducer2   = "test_producer2"
	tcConsumer1   = "test_consumer1"
	tcConsumer2   = "test_consumer2"
	tcWorkSize    = 22
	tcCapMessages = 1000
)

func TestError(t *testing.T) {
	str := "value"
	err := &SConfigError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func testConfigDefaultInit(configPath string) {
	_, _ = BuildConfig(configPath, &SConfig{
		FSettings: &SConfigSettings{
			FMessagesCapacity: tcCapMessages,
			FWorkSizeBits:     tcWorkSize,
			FNetworkKey:       tcNetwork,
		},
		FLogging: []string{"info", "erro"},
		FAddress: &SAddress{
			FHTTP:  tcAddress1,
			FPPROF: tcAddress2,
		},
		FProducers: []string{
			tcProducer1,
			tcProducer2,
		},
		FConsumers: []string{
			tcConsumer1,
			tcConsumer2,
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

	if cfg.GetSettings().GetMessagesCapacity() != tcCapMessages {
		t.Error("settings messages capacity is invalid")
		return
	}

	if cfg.GetSettings().GetWorkSizeBits() != tcWorkSize {
		t.Error("settings work size is invalid")
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

	for i, e := range []string{tcConsumer1, tcConsumer2} {
		if cfg.GetConsumers()[i] != e {
			t.Error("consumers is invalid")
			return
		}
	}

	for i, e := range []string{tcProducer1, tcProducer2} {
		if cfg.GetProducers()[i] != e {
			t.Error("producers is invalid")
			return
		}
	}
}
