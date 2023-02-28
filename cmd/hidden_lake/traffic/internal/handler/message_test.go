package handler

import (
	"fmt"
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload"
	testutils "github.com/number571/go-peer/test/_data"
)

func TestHandleMessageAPI(t *testing.T) {
	addr := testutils.TgAddrs[20]
	os.RemoveAll(fmt.Sprintf(databaseTemplate, addr))

	srv, connKeeper, db, hltClient := testAllRun(addr, "")
	defer testAllFree(addr, srv, connKeeper, db)

	client := testNewClient()
	msg, err := client.EncryptPayload(
		client.GetPubKey(),
		payload.NewPayload(0, []byte(testutils.TcLargeBody)),
	)
	if err != nil {
		t.Error(err)
		return
	}

	if err := hltClient.PutMessage(msg); err != nil {
		t.Error(err)
		return
	}

	strHash := encoding.HexEncode(msg.GetBody().GetHash())
	gotEncMsg, err := hltClient.GetMessage(strHash)
	if err != nil {
		t.Error(err)
		return
	}

	gotPubKey, gotPld, err := client.DecryptMessage(gotEncMsg)
	if err != nil {
		t.Error(err)
		return
	}

	if gotPubKey.Address().ToString() != client.GetPubKey().Address().ToString() {
		t.Error(err)
		return
	}

	if string(gotPld.GetBody()) != testutils.TcLargeBody {
		t.Error(err)
		return
	}
}
