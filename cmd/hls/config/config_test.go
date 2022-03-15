package config

import "testing"

const (
	tcConfigFile = "config_test.txt"
	tcAddress    = "test_address"
)

var (
	tgConnects = []string{
		"test_connect1",
		"test_connect2",
		"test_connect3",
	}
	tgServices = map[string]string{
		"test_service1": "test_address1",
		"test_service2": "test_address2",
		"test_service3": "test_address3",
	}
)

func TestConfig(t *testing.T) {
	cfg := NewConfig(tcConfigFile)

	if cfg.Address() != tcAddress {
		t.Errorf("address is invalid")
	}

	for i, v := range cfg.Connections() {
		if v != tgConnects[i] {
			t.Errorf("connection '%d' is invalid", i)
		}
	}

	for k, v := range tgServices {
		v1, ok := cfg.GetService(k)
		if !ok {
			t.Errorf("service undefined '%s'", k)
		}
		if v != v1 {
			t.Errorf("service is invalid '%s'", v1)
		}
	}
}
