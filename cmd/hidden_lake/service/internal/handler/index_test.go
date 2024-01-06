package handler

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	testutils "github.com/number571/go-peer/test/utils"
)

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
			fmt.Sprintf("http://%s", addr),
			&http.Client{Timeout: time.Minute},
		),
	)

	title, err := client.GetIndex()
	if err != nil {
		t.Error(err)
		return
	}

	if title != pkg_settings.CTitlePattern {
		t.Error("incorrect title pattern")
		return
	}
}
