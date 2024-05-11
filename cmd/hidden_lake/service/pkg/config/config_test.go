package config

import "testing"

func TestConfig(t *testing.T) {
	v := uint64(10)
	c := &SConfigSettings{
		FLimitMessageSizeBytes: v,
	}
	if c.GetLimitMessageSizeBytes() != v {
		t.Error("limit message size bytes != v")
		return
	}
}
