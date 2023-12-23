package handler

import (
	"fmt"
	"os"
	"testing"

	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	testutils "github.com/number571/go-peer/test/_data"
)

func TestHandleConfigSettingsAPI(t *testing.T) {
	t.Parallel()

	addr := testutils.TgAddrs[50]
	os.RemoveAll(fmt.Sprintf(databaseTemplate, addr))

	srv, connKeeper, cancel, db, hltClient := testAllRun(addr, "")
	defer testAllFree(addr, srv, connKeeper, cancel, db)

	settings, err := hltClient.GetSettings()
	if err != nil {
		t.Error(err)
		return
	}

	if settings.GetNetworkKey() != testutils.TCNetworkKey {
		t.Error("incorrect network key")
		return
	}

	if settings.GetKeySizeBits() != testutils.TcKeySize {
		t.Error("incorrect key size bits")
		return
	}

	if settings.GetWorkSizeBits() != testutils.TCWorkSize {
		t.Error("incorrect work size bits")
		return
	}

	if settings.GetMessagesCapacity() != testutils.TCCapacity {
		t.Error("incorrect messages capacity")
		return
	}

	if settings.GetMessageSizeBytes() != testutils.TCMessageSize {
		t.Error("incorrect messages size bytes")
		return
	}

	if settings.GetQueuePeriodMS() != hls_settings.CDefaultQueuePeriod {
		t.Error("incorrect queue period ms")
		return
	}

	if settings.GetLimitVoidSizeBytes() != hls_settings.CDefaultLimitVoidSize {
		t.Error("incorrect limit void size bytes")
		return
	}
}
