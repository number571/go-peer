package pprof

import "testing"

func TestPprof(t *testing.T) {
	t.Parallel()

	server := InitPprofService(":8080")
	if server.Addr != ":8080" {
		t.Error("incorrect address")
		return
	}
}
