package config

import (
	"fmt"
	"os"
	"testing"

	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
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

	if len(config1.GetServices()) != 3 {
		t.Error("got len(services) != 3")
		return
	}

	os.Remove(configFile)
	if err := os.WriteFile(configFile, []byte("abc"), 0o644); err != nil {
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

	config2, err := InitConfig(configFile, config1.(*SConfig))
	if err != nil {
		t.Error(err)
		return
	}

	if config2.GetServices()[0] != tgServices[0] {
		t.Error("got invalid field with exist config (2)")
		return
	}

	os.Remove(configFile)

	config3, err := InitConfig(configFile, nil)
	if err != nil {
		t.Error(err)
		return
	}

	if config3.GetServices()[0] != hls_settings.CTitlePattern {
		t.Error("got invalid field with exist config (3)")
		return
	}
}
