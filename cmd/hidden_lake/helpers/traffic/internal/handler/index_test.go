package handler

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/pkg/client"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/pkg/settings"
	"github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	testutils "github.com/number571/go-peer/test/utils"
)

func TestErrorsAPI(t *testing.T) {
	t.Parallel()

	client := hlt_client.NewClient(
		hlt_client.NewBuilder(),
		hlt_client.NewRequester(
			"http://"+testutils.TcUnknownHost,
			&http.Client{Timeout: time.Second},
			testNetworkMessageSettings(),
		),
	)

	pld := payload.NewPayload(uint64(testutils.TcHead), []byte(testutils.TcBody))
	sett := message.NewSettings(&message.SSettings{
		FWorkSizeBits: testutils.TCWorkSize,
	})
	if err := client.PutMessage(message.NewMessage(sett, pld, 1)); err == nil {
		t.Error("success put message with unknown host")
		return
	}

	if _, err := client.GetIndex(); err == nil {
		t.Error("success get index with unknown host")
		return
	}

	if _, err := client.GetHash(0); err == nil {
		t.Error("success get hash with unknown host")
		return
	}

	if _, err := client.GetMessage(""); err == nil {
		t.Error("success get message with unknown host")
		return
	}

	if _, err := client.GetPointer(); err == nil {
		t.Error("success get pointer with unknown host")
		return
	}

	if _, err := client.GetSettings(); err == nil {
		t.Error("success get settings with unknown host")
		return
	}
}

func TestHandleIndexAPI(t *testing.T) {
	t.Parallel()

	addr := testutils.TgAddrs[21]
	os.RemoveAll(fmt.Sprintf(databaseTemplate, addr))

	srv, connKeeper, cancel, db, hltClient := testAllRun(addr, "")
	defer testAllFree(addr, srv, connKeeper, cancel, db)

	title, err := hltClient.GetIndex()
	if err != nil {
		t.Error(err)
		return
	}

	if title != pkg_settings.CTitlePattern {
		t.Error("incorrect title pattern")
		return
	}
}
