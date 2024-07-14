package config

import (
	"os"
	"testing"
)

const (
	tcExecTimeout = 5000
	tcConfigFile  = "config_test.txt"
	tcLogging     = true
	tcAddress1    = "test_address1"
	tcAddress2    = "test_address2"
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
		FSettings: &SConfigSettings{
			FExecTimeoutMS: tcExecTimeout,
		},
		FLogging: []string{"info", "erro"},
		FAddress: &SAddress{
			FIncoming: tcAddress1,
			FPPROF:    tcAddress2,
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

	if cfg.GetAddress().GetIncoming() != tcAddress1 {
		t.Error("address incoming is invalid")
		return
	}

	if cfg.GetAddress().GetPPROF() != tcAddress2 {
		t.Error("address pprof is invalid")
		return
	}

	if cfg.GetSettings().GetExecTimeoutMS() != tcExecTimeout {
		t.Error("settings.exec_timeout_ms is invalid")
		return
	}
}
