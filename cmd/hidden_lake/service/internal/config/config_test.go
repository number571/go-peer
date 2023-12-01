package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	testutils "github.com/number571/go-peer/test/_data"
)

const (
	tcConfigFileTemplate = "config_test_%d.txt"
)

const (
	tcLogging       = true
	tcNetwork       = "test_network_key"
	tcDownloader    = "test_downloader"
	tcUploader      = "test_uploader"
	tcAddressTCP    = "test_address_tcp"
	tcAddressHTTP   = "test_address_http"
	tcAddressPPROF  = "test_address_pprof"
	tcPubKeyAlias1  = "test_alias1"
	tcPubKeyAlias2  = "test_alias2"
	tcServiceName1  = "test_service1"
	tcServiceName2  = "test_service2"
	tcMessageSize   = (1 << 20)
	tcWorkSize      = 20
	tcQueuePeriod   = 1000
	tcLimitVoidSize = (1 << 20)
)

var (
	tgConnects = []string{
		"test_connect1",
		"test_connect2",
	}
	tgPubKeys = map[string]string{
		tcPubKeyAlias1: testutils.TgPubKeys[0],
		tcPubKeyAlias2: testutils.TgPubKeys[1],
	}
	tgServices = map[string]string{
		tcServiceName1: "test_address1",
		tcServiceName2: "test_address2",
	}
)

const (
	tcConfigTemplate = `{
	"settings": {
		"message_size_bytes": %d,
		"work_size_bits": %d,
		"key_size_bits": %d,
		"queue_period_ms": %d,
		"limit_void_size_bytes": %d,
		"network_key": "%s"
	},
	"logging": ["info", "erro"],
	"address": {
		"tcp": "%s",
		"http": "%s",
		"pprof": "%s"
	},
	"connections": [
		"%s",
		"%s"
	],
	"friends": {
		"%s": "%s",
		"%s": "%s"
	},
	"services": {
		"%s": "%s",
		"%s": "%s"
	}
}`
)

func testNewConfigString() string {
	return fmt.Sprintf(
		tcConfigTemplate,
		tcMessageSize,
		tcWorkSize,
		testutils.TcKeySize,
		tcQueuePeriod,
		tcLimitVoidSize,
		tcNetwork,
		tcAddressTCP,
		tcAddressHTTP,
		tcAddressPPROF,
		tgConnects[0],
		tgConnects[1],
		tcPubKeyAlias1,
		tgPubKeys[tcPubKeyAlias1],
		tcPubKeyAlias2,
		tgPubKeys[tcPubKeyAlias2],
		tcServiceName1,
		tgServices[tcServiceName1],
		tcServiceName2,
		tgServices[tcServiceName2],
	)
}

func testConfigDefaultInit(configPath string) {
	os.WriteFile(configPath, []byte(testNewConfigString()), 0o644)
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

	if err := os.WriteFile(configFile, []byte("abc"), 0o644); err != nil {
		return err
	}

	if _, err := LoadConfig(configFile); err == nil {
		return errors.New("success load config with invalid structure")
	}

	cfg1Bytes := []byte(strings.ReplaceAll(testNewConfigString(), "settings", "settings_v2"))
	if err := os.WriteFile(configFile, cfg1Bytes, 0o644); err != nil {
		return err
	}

	if _, err := LoadConfig(configFile); err == nil {
		return errors.New("success load config with required fields (settings)")
	}

	cfg2Bytes := []byte(strings.ReplaceAll(testNewConfigString(), "PubKey", "PubKey_v2"))
	if err := os.WriteFile(configFile, cfg2Bytes, 0o644); err != nil {
		return err
	}

	if _, err := LoadConfig(configFile); err == nil {
		return errors.New("success load config with invalid fields (friends)")
	}

	cfg3Bytes := []byte(strings.ReplaceAll(testNewConfigString(), "erro", "erro_v2"))
	if err := os.WriteFile(configFile, cfg3Bytes, 0o644); err != nil {
		return err
	}

	if _, err := LoadConfig(configFile); err == nil {
		return errors.New("success load config with invalid fields (logging)")
	}

	pubKey1 := tgPubKeys[tcPubKeyAlias1]
	pubKey2 := tgPubKeys[tcPubKeyAlias2]

	cfg4Bytes := []byte(strings.ReplaceAll(testNewConfigString(), pubKey1, pubKey2))
	if err := os.WriteFile(configFile, cfg4Bytes, 0o644); err != nil {
		return err
	}

	if _, err := LoadConfig(configFile); err == nil {
		return errors.New("success load config with invalid fields (duplicate publc keys)")
	}

	newPubKey := asymmetric.NewRSAPrivKey(512).GetPubKey().ToString()
	cfg5Bytes := []byte(strings.ReplaceAll(testNewConfigString(), pubKey1, newPubKey))
	if err := os.WriteFile(configFile, cfg5Bytes, 0o644); err != nil {
		return err
	}

	if _, err := LoadConfig(configFile); err == nil {
		return errors.New("success load config with invalid fields (diff key sizes)")
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

	if cfg.GetSettings().GetKeySizeBits() != testutils.TcKeySize {
		t.Error("settings key size is invalid")
		return
	}

	if cfg.GetSettings().GetQueuePeriodMS() != tcQueuePeriod {
		t.Error("settings queue period is invalid")
		return
	}

	if cfg.GetSettings().GetLimitVoidSizeBytes() != tcLimitVoidSize {
		t.Error("settings limit void size is invalid")
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
		if v != v1 {
			t.Errorf("service address is invalid '%s'", v1)
			return
		}
	}

	for name, pubStr := range tgPubKeys {
		v1 := cfg.GetFriends()[name]
		pubKey := asymmetric.LoadRSAPubKey(pubStr)
		if pubKey.GetAddress().ToString() != v1.GetAddress().ToString() {
			t.Errorf("public key is invalid '%s'", v1)
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
	wrapper.GetEditor().UpdateFriends(nil)

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
