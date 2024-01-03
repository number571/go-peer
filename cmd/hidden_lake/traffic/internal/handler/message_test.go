package handler

import (
	"fmt"
	"os"
	"testing"

	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/encoding"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	testutils "github.com/number571/go-peer/test/_data"
)

func TestHandleMessageAPI(t *testing.T) {
	t.Parallel()

	addr := testutils.TgAddrs[20]
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

	strHash := encoding.HexEncode(netMsg.GetHash())
	gotNetMsg, err := hltClient.GetMessage(strHash)
	if err != nil {
		t.Error(err)
		return
	}

	gotEncMsg, err := message.LoadMessage(
		client.GetSettings(),
		gotNetMsg.GetPayload().GetBody(),
	)
	if err != nil {
		t.Error(err)
		return
	}

	gotPubKey, gotPld, err := client.DecryptMessage(gotEncMsg)
	if err != nil {
		t.Error(err)
		return
	}

	if gotPubKey.GetAddress().ToString() != client.GetPubKey().GetAddress().ToString() {
		t.Error(err)
		return
	}

	if string(gotPld.GetBody()) != testutils.TcBody {
		t.Error(err)
		return
	}
}
