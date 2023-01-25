package handler

import (
	"testing"

	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/settings"
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload"
	testutils "github.com/number571/go-peer/test/_data"
)

func TestHandleMessageAPI(t *testing.T) {
	addr := testutils.TgAddrs[20]

	srv, db, hltClient := testAllRun(addr)
	defer testAllFree(addr, srv, db)

	privKey := asymmetric.LoadRSAPrivKey(testutils.TcPrivKey)
	pubKey := privKey.PubKey()

	client := client.NewClient(client.NewSettings(
		&client.SSettings{
			FMessageSize: hlt_settings.CMessageSize,
			FWorkSize:    hlt_settings.CWorkSize,
		}),
		privKey,
	)

	msg, err := client.Encrypt(
		pubKey,
		payload.NewPayload(0, []byte(testutils.TcLargeBody)),
	)
	if err != nil {
		t.Error(err)
		return
	}

	if err := hltClient.AddMessage(msg); err != nil {
		t.Error(err)
		return
	}

	strHash := encoding.HexEncode(msg.Body().Hash())
	gotEncMsg, err := hltClient.GetMessage(strHash)
	if err != nil {
		t.Error(err)
		return
	}

	gotPubKey, gotPld, err := client.Decrypt(gotEncMsg)
	if err != nil {
		t.Error(err)
		return
	}

	if gotPubKey.Address().String() != pubKey.Address().String() {
		t.Error(err)
		return
	}

	if string(gotPld.Body()) != testutils.TcLargeBody {
		t.Error(err)
		return
	}
}