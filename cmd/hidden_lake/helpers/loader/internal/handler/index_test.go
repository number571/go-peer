package handler

import (
	"net/http"
	"testing"
	"time"

	hll_client "github.com/number571/go-peer/cmd/hidden_lake/helpers/loader/pkg/client"
	"github.com/number571/go-peer/cmd/hidden_lake/helpers/loader/pkg/settings"
	testutils "github.com/number571/go-peer/test/utils"
)

func TestErrorsAPI(t *testing.T) {
	t.Parallel()

	client := hll_client.NewClient(
		hll_client.NewRequester(
			"http://"+testutils.TcUnknownHost,
			&http.Client{Timeout: time.Second},
		),
	)

	if err := client.RunTransfer(); err == nil {
		t.Error("success run transfer with unknown host")
		return
	}

	if err := client.StopTransfer(); err == nil {
		t.Error("success stop transfer with unknown host")
		return
	}

	if _, err := client.GetIndex(); err == nil {
		t.Error("success get index with unknown host")
		return
	}

	if _, err := client.GetSettings(); err == nil {
		t.Error("success get settings with unknown host")
		return
	}
}

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
