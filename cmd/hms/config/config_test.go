package config

import "testing"

const (
	tcConfigFile = "config_test.txt"
	tcAddress    = "test_address"
	tcCleanCron  = "0 0 0 0 0"
)

func TestConfig(t *testing.T) {
	cfg := NewConfig(tcConfigFile)

	if cfg.Address() != tcAddress {
		t.Errorf("address is invalid")
	}

	if cfg.CleanCron() != tcCleanCron {
		t.Errorf("clean_cron is invalid")
	}
}
