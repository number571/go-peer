package config

import (
	"os"
	"testing"
)

const (
	tcConfigFile      = "config_test.txt"
	tcLogging         = true
	tcAddress1        = "test_address1"
	tcAddress2        = "test_address2"
	tcService1        = "test_service1"
	tcServiceAddress1 = "test_service_address1"
	tcWorkSize        = 22
	tcCapMessages     = 1000
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SConfigError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func testConfigDefaultInit(configPath string) {
	_, _ = BuildConfig(configPath, &SConfig{
		FSettings: &SConfigSettings{},
		FLogging:  []string{"info", "erro"},
		FAddress: &SAddress{
			FHTTP:  tcAddress1,
			FPPROF: tcAddress2,
		},
		FServices: map[string]string{
			tcService1: tcServiceAddress1,
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

	if cfg.GetAddress().GetHTTP() != tcAddress1 {
		t.Error("address http is invalid")
		return
	}

	if cfg.GetAddress().GetPPROF() != tcAddress2 {
		t.Error("address pprof is invalid")
		return
	}

	if v, ok := cfg.GetService(tcService1); !ok || v != tcServiceAddress1 {
		t.Error("service address is invalid")
		return
	}
}
