package handler

import (
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

	srv, connKeeper, cancel, db, hltClient := testAllRun(addr, "")
	defer testAllFree(addr, srv, connKeeper, cancel, db)

	client := testNewClient()
	msg, err := client.EncryptPayload(
		client.GetPubKey(),
		payload.NewPayload(0, []byte(testutils.TcBody)),
	)
	if err != nil {
		t.Error(err)
		return
	}

	netMsg := net_message.NewMessage(
		testNetworkMessageSettings(),
		payload.NewPayload(hls_settings.CNetworkMask, msg.ToBytes()),
		1,
	)
	if err := hltClient.PutMessage(netMsg); err != nil {
		t.Error(err)
		return
	}

	pointer, err := hltClient.GetPointer()
	if err != nil {
		t.Error(err)
		return
	}

	if pointer != 1 {
		t.Error("incorrect pointer")
		return
	}
}
