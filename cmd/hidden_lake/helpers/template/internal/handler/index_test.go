package handler

import (
	"context"
	"net/http"
	"testing"
	"time"

	hl_t_client "github.com/number571/go-peer/cmd/hidden_lake/helpers/template/pkg/client"
	hl_t_settings "github.com/number571/go-peer/cmd/hidden_lake/helpers/template/pkg/settings"
	testutils "github.com/number571/go-peer/test/utils"
)

func TestErrorsAPI(t *testing.T) {
	t.Parallel()

	client := hl_t_client.NewClient(
		hl_t_client.NewRequester(
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

	service := testRunService(testutils.TgAddrs[53])
	defer service.Close()

	time.Sleep(100 * time.Millisecond)
	hleClient := hl_t_client.NewClient(
		hl_t_client.NewRequester(
			"http://"+testutils.TgAddrs[53],
			&http.Client{Timeout: time.Second / 2},
		),
	)

	title, err := hleClient.GetIndex(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	if title != hl_t_settings.CServiceFullName {
		t.Error("incorrect title pattern")
		return
	}
}
