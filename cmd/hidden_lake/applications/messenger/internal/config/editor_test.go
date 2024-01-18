package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/number571/go-peer/internal/language"
	logger "github.com/number571/go-peer/internal/logger/std"
	"github.com/number571/go-peer/pkg/crypto/random"
)

const (
	tcConfigFileTemplate = "config_test_%d.txt"
)

type tsConfig struct{}

var (
	_ IConfig = &tsConfig{}
)

func (p *tsConfig) GetSettings() IConfigSettings     { return nil }
func (p *tsConfig) GetLanguage() language.ILanguage  { return 0 }
func (p *tsConfig) GetLogging() logger.ILogging      { return nil }
func (p *tsConfig) GetShare() bool                   { return false }
func (p *tsConfig) GetPseudonym() string             { return "" }
func (p *tsConfig) GetAddress() IAddress             { return nil }
func (p *tsConfig) GetNetworkKey() string            { return "" }
func (p *tsConfig) GetConnection() string            { return "" }
func (p *tsConfig) GetStorageKey() string            { return "" }
func (p *tsConfig) GetSecretKeys() map[string]string { return nil }

func TestPanicEditor(t *testing.T) {
	t.Parallel()

	for i := 0; i < 2; i++ {
		testPanicEditor(t, i)
	}
}

func testPanicEditor(t *testing.T, n int) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()
	switch n {
	case 0:
		_ = newEditor(nil)
	case 1:
		_ = newEditor(&tsConfig{})
	}
}

func TestIncorrectFilepathEditor(t *testing.T) {
	t.Parallel()

	configFile := fmt.Sprintf(tcConfigFileTemplate, 3)
	defer os.Remove(configFile)

	testConfigDefaultInit(configFile)
	cfg, err := LoadConfig(configFile)
	if err != nil {
		t.Error(err)
		return
	}

	wrapper := NewWrapper(cfg)

	config := wrapper.GetConfig().(*SConfig)
	editor := wrapper.GetEditor()

	config.fFilepath = random.NewStdPRNG().GetString(32)

	res, err := language.ToILanguage("RUS")
	if err != nil {
		t.Error(err)
		return
	}
	if err := editor.UpdateLanguage(res); err == nil {
		t.Error("success update network key with incorrect filepath")
		return
	}

	if err := editor.UpdatePseudonym("test"); err == nil {
		t.Error("success update pseudonym with incorrect filepath")
		return
	}

	if err := editor.UpdateSecretKeys(map[string]string{"Alice": "123"}); err == nil {
		t.Error("success update friends with incorrect filepath")
		return
	}
}
