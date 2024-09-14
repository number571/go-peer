package handler

import (
	"context"
	"net/http"
	"testing"
	"time"

	hld_client "github.com/number571/go-peer/cmd/hidden_lake/helpers/distributor/pkg/client"
	testutils "github.com/number571/go-peer/test/utils"
)

func TestHandleConfigSettingsAPI(t *testing.T) {
	t.Parallel()

	service := testRunService(testutils.TgAddrs[53], "")
	defer service.Close()

	time.Sleep(100 * time.Millisecond)
	hldClient := hld_client.NewClient(
		hld_client.NewRequester(
			"http://"+testutils.TgAddrs[53],
			&http.Client{Timeout: time.Second / 2},
		),
	)

	settings, err := hldClient.GetSettings(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	_ = settings
}
