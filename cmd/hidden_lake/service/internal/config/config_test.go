// nolint: goerr113
package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/number571/go-peer/pkg/encoding"
)

const (
	tcConfigFileTemplate = "config_test_%d.txt"
)

const (
	tcLogging         = true
	tcNetwork         = "test_network_key"
	tcDownloader      = "test_downloader"
	tcUploader        = "test_uploader"
	tcAddressTCP      = "test_address_tcp"
	tcAddressHTTP     = "test_address_http"
	tcAddressPPROF    = "test_address_pprof"
	tcKeyAlias1       = "test_alias1"
	tcKeyAlias2       = "test_alias2"
	tcServiceName1    = "test_service1"
	tcServiceName2    = "test_service2"
	tcMessageSize     = (1 << 20)
	tcWorkSize        = 22
	tcFetchTimeout    = 5000
	tcQueuePeriod     = 1000
	tcQueueRandPeriod = 2000
	tcLimitVoidSize   = (1 << 20)
)

var (
	tgConnects = []string{
		"test_connect1",
		"test_connect2",
	}
	tgKeys = map[string]string{
		tcKeyAlias1: encoding.HexEncode([]byte("helloworld1672y8hdhf4328eh191d21")),
		tcKeyAlias2: encoding.HexEncode([]byte("f4uhf278hd2u3d8u2fh438fhj39dj2ii")),
	}
	tgServices = map[string]string{
		tcServiceName1: "test_address1",
		tcServiceName2: "test_address2",
	}
)

const (
	tcConfigTemplate = `settings:
  message_size_bytes: %d
  work_size_bits: %d
  fetch_timeout_ms: %d
  queue_period_ms: %d
  rand_queue_period_ms: %d
  rand_message_size_bytes: %d
  network_key: %s
  f2f_disabled: true
logging:
  - info
  - erro
address:
  tcp: %s
  http: %s
  pprof: %s
connections:
  - %s
  - %s
friends:
  %s: %s
  %s: %s
services:
  %s: 
    host: %s
  %s: 
    host: %s
`
)

func testNewConfigString() string {
	return fmt.Sprintf(
		tcConfigTemplate,
		tcMessageSize,
		tcWorkSize,
		tcFetchTimeout,
		tcQueuePeriod,
		tcQueueRandPeriod,
		tcLimitVoidSize,
		tcNetwork,
		tcAddressTCP,
		tcAddressHTTP,
		tcAddressPPROF,
		tgConnects[0],
		tgConnects[1],
		tcKeyAlias1,
		tgKeys[tcKeyAlias1],
		tcKeyAlias2,
		tgKeys[tcKeyAlias2],
		tcServiceName1,
		tgServices[tcServiceName1],
		tcServiceName2,
		tgServices[tcServiceName2],
	)
}

func testConfigDefaultInit(configPath string) {
	_ = os.WriteFile(configPath, []byte(testNewConfigString()), 0o600)
}

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SConfigError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestBuildConfig(t *testing.T) {
	t.Parallel()

	config1File := fmt.Sprintf(tcConfigFileTemplate, 2)
	config2File := fmt.Sprintf(tcConfigFileTemplate, 3)

	testConfigDefaultInit(config1File)
	defer os.Remove(config1File)

	cfg, err := LoadConfig(config1File)
	if err != nil {
		t.Error(err)
		return
	}

	if _, err := BuildConfig(config2File, &SConfig{}); err == nil {
		t.Error("success build config with void structure")
		return
	}

	if _, err := BuildConfig(config2File, cfg.(*SConfig)); err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(config2File)

	if _, err := BuildConfig(config2File, cfg.(*SConfig)); err == nil {
		t.Error("success build already exist config")
		return
	}
}

func testIncorrectConfig(configFile string) error {
	if _, err := LoadConfig(configFile); err == nil {
		return errors.New("success load config on non exist file")
	}

	if err := os.WriteFile(configFile, []byte("abc"), 0o600); err != nil {
		return err
	}

	if _, err := LoadConfig(configFile); err == nil {
		return errors.New("success load config with invalid structure")
	}

	cfg1Bytes := []byte(strings.ReplaceAll(testNewConfigString(), "settings", "settings_v2"))
	if err := os.WriteFile(configFile, cfg1Bytes, 0o600); err != nil {
		return err
	}

	if _, err := LoadConfig(configFile); err == nil {
		return errors.New("success load config with required fields (settings)")
	}

	cfg3Bytes := []byte(strings.ReplaceAll(testNewConfigString(), "erro", "erro_v2"))
	if err := os.WriteFile(configFile, cfg3Bytes, 0o600); err != nil {
		return err
	}

	if _, err := LoadConfig(configFile); err == nil {
		return errors.New("success load config with invalid fields (logging)")
	}

	pubKey1 := tgKeys[tcKeyAlias1]
	pubKey2 := tgKeys[tcKeyAlias2]

	cfg4Bytes := []byte(strings.ReplaceAll(testNewConfigString(), pubKey1, pubKey2))
	if err := os.WriteFile(configFile, cfg4Bytes, 0o600); err != nil {
		return err
	}

	if _, err := LoadConfig(configFile); err == nil {
		return errors.New("success load config with invalid fields (duplicate publc keys)")
	}

	return nil
}

func TestComplexConfig(t *testing.T) {
	t.Parallel()

	configFile := fmt.Sprintf(tcConfigFileTemplate, 0)
	defer os.Remove(configFile)

	if err := testIncorrectConfig(configFile); err != nil {
		t.Error(err)
		return
	}

	testConfigDefaultInit(configFile)
	cfg, err := LoadConfig(configFile)
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

	if cfg.GetSettings().GetFetchTimeoutMS() != tcFetchTimeout {
		t.Error("settings fetch timeout is invalid")
		return
	}

	if cfg.GetSettings().GetQueuePeriodMS() != tcQueuePeriod {
		t.Error("settings queue period is invalid")
		return
	}

	if cfg.GetSettings().GetRandQueuePeriodMS() != tcQueueRandPeriod {
		t.Error("settings rand queue period is invalid")
		return
	}

	if cfg.GetSettings().GetRandMessageSizeBytes() != tcLimitVoidSize {
		t.Error("settings rand message size is invalid")
		return
	}

	if !cfg.GetSettings().GetF2FDisabled() {
		t.Error("settings f2f disabled is invalid")
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

	if cfg.GetAddress().GetTCP() != tcAddressTCP {
		t.Error("address_tcp is invalid")
		return
	}

	if cfg.GetAddress().GetHTTP() != tcAddressHTTP {
		t.Error("address_http is invalid")
		return
	}

	if cfg.GetAddress().GetPPROF() != tcAddressPPROF {
		t.Error("address_pprof is invalid")
		return
	}

	if len(cfg.GetConnections()) != 2 {
		t.Error("len connections != 2")
		return
	}
	for i, v := range cfg.GetConnections() {
		if v != tgConnects[i] {
			t.Errorf("connection '%d' is invalid", i)
			return
		}
	}

	for k, v := range tgServices {
		v1, ok := cfg.GetService(k)
		if !ok {
			t.Errorf("service undefined '%s'", k)
			return
		}
		if v != v1.GetHost() {
			t.Errorf("service host is invalid '%s'", v1)
			return
		}
	}

	for name, keyStr := range tgKeys {
		v1 := cfg.GetFriends()[name]
		if keyStr != v1 {
			t.Errorf("key is invalid '%s'", v1)
			return
		}
	}
}

func TestWrapper(t *testing.T) {
	t.Parallel()

	configFile := fmt.Sprintf(tcConfigFileTemplate, 1)

	testConfigDefaultInit(configFile)
	defer os.Remove(configFile)

	cfg, err := LoadConfig(configFile)
	if err != nil {
		t.Error(err)
		return
	}

	if len(cfg.GetFriends()) == 0 {
		t.Error("list of friends should be is not nil for tests")
		return
	}

	wrapper := NewWrapper(cfg)
	_ = wrapper.GetEditor().UpdateFriends(nil)

	if len(cfg.GetFriends()) != 0 {
		t.Error("friends is not nil for current config")
		return
	}

	cfg, err = LoadConfig(configFile)
	if err != nil {
		t.Error(err)
		return
	}

	if len(cfg.GetFriends()) != 0 {
		t.Error("friends is not nil for loaded config")
		return
	}
}
