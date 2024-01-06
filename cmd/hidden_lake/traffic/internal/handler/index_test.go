package handler

import (
	"fmt"
	"os"
	"testing"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	testutils "github.com/number571/go-peer/test/utils"
)

func TestHandleIndexAPI(t *testing.T) {
	t.Parallel()

	addr := testutils.TgAddrs[21]
	os.RemoveAll(fmt.Sprintf(databaseTemplate, addr))

	srv, connKeeper, cancel, db, hltClient := testAllRun(addr, "")
	defer testAllFree(addr, srv, connKeeper, cancel, db)

	title, err := hltClient.GetIndex()
	if err != nil {
		t.Error(err)
		return
	}

	if title != pkg_settings.CTitlePattern {
		t.Error("incorrect title pattern")
		return
	}
}
