package handler

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	testutils "github.com/number571/go-peer/test/utils"
)

func TestErrorsAPI(t *testing.T) {
	t.Parallel()

	client := hls_client.NewClient(
		hls_client.NewBuilder(),
		hls_client.NewRequester(
			"http://"+testutils.TcUnknownHost,
			&http.Client{Timeout: time.Second},
		),
	)

	if err := client.AddConnection(context.Background(), ""); err == nil {
		t.Error("success add connection with unknown host")
		return
	}

	if err := client.DelConnection(context.Background(), ""); err == nil {
		t.Error("success del connection with unknown host")
		return
	}

	if err := client.AddFriend(context.Background(), "", asymmetric.LoadRSAPubKey(testutils.TgPubKeys[0])); err == nil {
		t.Error("success add friend with unknown host")
		return
	}

	if err := client.DelFriend(context.Background(), ""); err == nil {
		t.Error("success del friend with unknown host")
		return
	}

	if err := client.BroadcastRequest(context.Background(), "", request.NewRequest("", "", "")); err == nil {
		t.Error("success broadcast request with unknown host")
		return
	}

	if _, err := client.FetchRequest(context.Background(), "", request.NewRequest("", "", "")); err == nil {
		t.Error("success fetch request with unknown host")
		return
	}

	if _, err := client.GetIndex(context.Background()); err == nil {
		t.Error("success get index with unknown host")
		return
	}

	if _, err := client.GetConnections(context.Background()); err == nil {
		t.Error("success get connections with unknown host")
		return
	}

	if _, err := client.GetFriends(context.Background()); err == nil {
		t.Error("success get friends with unknown host")
		return
	}

	if _, err := client.GetOnlines(context.Background()); err == nil {
		t.Error("success get onlines with unknown host")
		return
	}

	if _, err := client.GetPubKey(context.Background()); err == nil {
		t.Error("success get pub key with unknown host")
		return
	}

	if _, err := client.GetSettings(context.Background()); err == nil {
		t.Error("success get settings with unknown host")
		return
	}

	if err := client.SetNetworkKey(context.Background(), "test"); err == nil {
		t.Error("success set network key with unknown host")
		return
	}

	if err := client.DelOnline(context.Background(), "test"); err == nil {
		t.Error("success del online key with unknown host")
		return
	}
}

func TestHandleIndexAPI(t *testing.T) {
	t.Parallel()

	addr := testutils.TgAddrs[22]
	pathCfg := fmt.Sprintf(tcPathConfigTemplate, 3)
	pathDB := fmt.Sprintf(tcPathDBTemplate, 3)

	_, node, _, cancel, srv := testAllCreate(pathCfg, pathDB, addr)
	defer testAllFree(node, cancel, srv, pathCfg, pathDB)

	client := hls_client.NewClient(
		hls_client.NewBuilder(),
		hls_client.NewRequester(
			"http://"+addr,
			&http.Client{Timeout: time.Minute},
		),
	)

	title, err := client.GetIndex(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	if title != pkg_settings.CServiceFullName {
		t.Error("incorrect title pattern")
		return
	}
}
