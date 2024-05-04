package handler

import (
	"context"
	"fmt"
	"os"
	"testing"

	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	testutils "github.com/number571/go-peer/test/utils"
)

func TestHandlePointerAPI(t *testing.T) {
	t.Parallel()

	addr := testutils.TgAddrs[46]
	os.RemoveAll(fmt.Sprintf(databaseTemplate, addr))

	srv, cancel, db, hltClient := testAllRun(addr)
	defer testAllFree(addr, srv, cancel, db)

	client := testNewClient()
	msg, err := client.EncryptPayload(
		client.GetPubKey(),
		payload.NewPayload64(0, []byte(testutils.TcBody)),
	)
	if err != nil {
		t.Error(err)
		return
	}

	netMsg := net_message.NewMessage(
		testNetworkMessageSettings(),
		payload.NewPayload64(hls_settings.CNetworkMask, msg.ToBytes()),
	)
	if err := hltClient.PutMessage(context.Background(), netMsg); err != nil {
		t.Error(err)
		return
	}

	pointer, err := hltClient.GetPointer(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	if pointer != 1 {
		t.Error("incorrect pointer")
		return
	}
}
