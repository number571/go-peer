package config

import (
	"fmt"
	"os"
	"testing"

	hlr_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/remoter/pkg/settings"
)

const (
	tcConfigFileTemplate = "config_test_%d.txt"
)

func TestInit(t *testing.T) {
	t.Parallel()

	configFile := fmt.Sprintf(tcConfigFileTemplate, 1)
	defer os.Remove(configFile)

	testConfigDefaultInit(configFile)

	config1, err := InitConfig(configFile, nil)
	if err != nil {
		t.Error(err)
		return
	}

	if config1.GetAddress().GetIncoming() != tcAddress1 {
		t.Error("got invalid field with exist config (1)")
		return
	}

	if config1.GetAddress().GetPPROF() != tcAddress2 {
		t.Error("got invalid field with exist config (2)")
		return
	}

	os.Remove(configFile)
	if err := os.WriteFile(configFile, []byte("abc"), 0o600); err != nil {
		t.Error(err)
		return
	}

	if _, err := InitConfig(configFile, nil); err == nil {
		t.Error("success init config with invalid config structure (1)")
		return
	}

	os.Remove(configFile)

	if _, err := InitConfig(configFile, &SConfig{}); err == nil {
		t.Error("success init config with invalid config structure (2)")
		return
	}

	os.Remove(configFile)

	config3, err := InitConfig(configFile, nil)
	if err != nil {
		t.Error(err)
		return
	}

	if config3.GetAddress().GetIncoming() != hlr_settings.CDefaultIncomingAddress {
		t.Error("got invalid field with exist config (3)")
		return
	}
}
