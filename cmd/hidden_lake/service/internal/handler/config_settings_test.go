package handler

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	testutils "github.com/number571/go-peer/test/utils"
)

func TestHandleConfigSettingsAPI(t *testing.T) {
	t.Parallel()

	addr := testutils.TgAddrs[26]
	pathCfg := fmt.Sprintf(tcPathConfigTemplate, 2)
	pathDB := fmt.Sprintf(tcPathDBTemplate, 2)

	_, node, _, cancel, srv := testAllCreate(pathCfg, pathDB, addr)
	defer testAllFree(node, cancel, srv, pathCfg, pathDB)

	client := hls_client.NewClient(
		hls_client.NewBuilder(),
		hls_client.NewRequester(
			fmt.Sprintf("http://%s", addr),
			&http.Client{Timeout: time.Minute},
		),
	)

	sett, err := client.GetSettings(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	if sett.GetKeySizeBits() != testutils.TcKeySize {
		t.Error("invalid key size")
		return
	}

	if sett.GetQueuePeriodMS() != 1000 {
		t.Error("invalid queue period")
		return
	}

	if sett.GetLimitVoidSizeBytes() != 4096 {
		t.Error("invalid limit void size")
		return
	}

	if sett.GetMessageSizeBytes() != 8192 {
		t.Error("invalid message size")
		return
	}

	if sett.GetWorkSizeBits() != 22 {
		t.Error("invalid work size")
		return
	}
}
