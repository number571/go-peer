package handler

import (
	"testing"

	testutils "github.com/number571/go-peer/test/_data"
)

func TestBroadcastIndexAPI(t *testing.T) {
	addr := testutils.TgAddrs[23]

	srv, db, hltClient := testAllRun(addr)
	defer testAllFree(addr, srv, db)

	if err := hltClient.DoBroadcast(); err != nil {
		t.Error(err)
		return
	}
}
