package handler

import (
	"net/http"
	"testing"
	"time"

	hl_t_client "github.com/number571/go-peer/cmd/hidden_lake/_template/pkg/client"
	testutils "github.com/number571/go-peer/test/_data"
)

func TestHandleConfigSettingsAPI(t *testing.T) {
	t.Parallel()

	service := testRunService(testutils.TgAddrs[51])
	defer service.Close()

	time.Sleep(100 * time.Millisecond)
	hleClient := hl_t_client.NewClient(
		hl_t_client.NewRequester(
			"http://"+testutils.TgAddrs[51],
			&http.Client{Timeout: time.Second / 2},
		),
	)

	settings, err := hleClient.GetSettings()
	if err != nil {
		t.Error(err)
		return
	}

	if settings.GetValue() != "TODO" {
		t.Error("incorrect value")
		return
	}
}