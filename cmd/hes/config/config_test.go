package config

import "testing"

const (
	tcConfigFile = "config_test.txt"
	tcAddress    = "test_address"
)

func TestConfig(t *testing.T) {
	cfg := NewConfig(tcConfigFile)

	if cfg.Address() != tcAddress {
		t.Errorf("address is invalid")
	}
}
