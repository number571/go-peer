package handler

import (
	"net/http"
	"testing"
	"time"

	hll_client "github.com/number571/go-peer/cmd/hidden_lake/helpers/loader/pkg/client"
	"github.com/number571/go-peer/cmd/hidden_lake/helpers/loader/pkg/settings"
	testutils "github.com/number571/go-peer/test/_data"
)

func TestHandleIndexAPI(t *testing.T) {
	t.Parallel()

	service := testRunService(testutils.TgAddrs[45])
	defer service.Close()

	time.Sleep(100 * time.Millisecond)
	hllClient := hll_client.NewClient(
		hll_client.NewRequester(
			"http://"+testutils.TgAddrs[45],
			&http.Client{Timeout: time.Second / 2},
		),
	)

	title, err := hllClient.GetIndex()
	if err != nil {
		t.Error(err)
		return
	}

	if title != settings.CTitlePattern {
		t.Error("incorrect title pattern")
		return
	}
}
