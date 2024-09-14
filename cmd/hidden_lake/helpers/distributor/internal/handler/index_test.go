package handler

import (
	"context"
	"net/http"
	"testing"
	"time"

	hld_client "github.com/number571/go-peer/cmd/hidden_lake/helpers/distributor/pkg/client"
	"github.com/number571/go-peer/cmd/hidden_lake/helpers/distributor/pkg/settings"
	testutils "github.com/number571/go-peer/test/utils"
)

func TestErrorsAPI(t *testing.T) {
	t.Parallel()

	client := hld_client.NewClient(
		hld_client.NewRequester(
			"http://"+testutils.TcUnknownHost,
			&http.Client{Timeout: time.Second},
		),
	)

	if _, err := client.GetIndex(context.Background()); err == nil {
		t.Error("success get index with unknown host")
		return
	}

	if _, err := client.GetSettings(context.Background()); err == nil {
		t.Error("success get settings with unknown host")
		return
	}
}

func TestHandleIndexAPI(t *testing.T) {
	t.Parallel()

	service := testRunService(testutils.TgAddrs[51], "")
	defer service.Close()

	time.Sleep(100 * time.Millisecond)
	hldClient := hld_client.NewClient(
		hld_client.NewRequester(
			"http://"+testutils.TgAddrs[51],
			&http.Client{Timeout: time.Second / 2},
		),
	)

	title, err := hldClient.GetIndex(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	if title != settings.CServiceFullName {
		t.Error("incorrect title pattern")
		return
	}
}
