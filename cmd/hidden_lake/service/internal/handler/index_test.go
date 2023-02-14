package handler

import (
	"fmt"
	"testing"

	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	testutils "github.com/number571/go-peer/test/_data"
)

func TestHandleIndexAPI(t *testing.T) {
	addr := testutils.TgAddrs[22]

	_, node, srv := testAllCreate(tcPathConfig, tcPathDB, addr)
	defer testAllFree(node, srv)

	client := hls_client.NewClient(
		hls_client.NewBuilder(),
		hls_client.NewRequester(fmt.Sprintf("http://%s", addr)),
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
