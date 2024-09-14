package handler

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	hld_client "github.com/number571/go-peer/cmd/hidden_lake/helpers/distributor/pkg/client"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	testutils "github.com/number571/go-peer/test/utils"
)

const (
	tcServiceHost      = "service-host"
	tcServicePath      = "/some-service-path"
	tcResponseMsg      = "(service-response-msg)"
	tcServiceHeadKey   = "Service-Header-Key"
	tcServiceHeadValue = "HeadValue"
	tcServiceBody      = "request-body"
)

func TestHandleNetworkDistributeAPI(t *testing.T) {
	t.Parallel()

	service := testRunService(testutils.TgAddrs[60], testutils.TgAddrs[61])
	defer service.Close()

	innerService := testRunInnerService(testutils.TgAddrs[61])
	defer innerService.Close()

	time.Sleep(100 * time.Millisecond)
	hldClient := hld_client.NewClient(
		hld_client.NewRequester(
			"http://"+testutils.TgAddrs[60],
			&http.Client{Timeout: time.Second / 2},
		),
	)

	resp, err := hldClient.DistributeRequest(
		context.Background(),
		request.NewRequest(http.MethodGet, tcServiceHost, tcServicePath).
			WithHead(map[string]string{
				tcServiceHeadKey: tcServiceHeadValue,
			}).
			WithBody([]byte(tcServiceBody)),
	)
	if err != nil {
		t.Error(err)
		return
	}

	if resp.GetCode() != http.StatusAccepted {
		t.Error("resp.GetCode() != http.StatusAccepted")
		return
	}

	if !bytes.Equal(resp.GetBody(), []byte(fmt.Sprintf("%s: %s", tcResponseMsg, tcServiceBody))) {
		t.Error("bad response body")
		return
	}
}
