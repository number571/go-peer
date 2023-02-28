package handler

import (
	"fmt"
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload"
	testutils "github.com/number571/go-peer/test/_data"
	anon_testutils "github.com/number571/go-peer/test/_data/anonymity"
)

func TestHandleHashesAPI(t *testing.T) {
	addr := testutils.TgAddrs[19]
	os.RemoveAll(fmt.Sprintf(databaseTemplate, addr))

	srv, connKeeper, db, hltClient := testAllRun(addr, "")
	defer testAllFree(addr, srv, connKeeper, db)

	privKey := asymmetric.LoadRSAPrivKey(testutils.TcPrivKey)
	pubKey := asymmetric.LoadRSAPubKey(testutils.TgPubKeys[0])

	client := client.NewClient(client.NewSettings(
		&client.SSettings{
			FMessageSize: anon_testutils.TCMessageSize,
			FWorkSize:    anon_testutils.TCWorkSize,
		}),
		privKey,
	)

	msg, err := client.EncryptPayload(
		pubKey,
		payload.NewPayload(0, []byte("hello")),
	)
	if err != nil {
		t.Error(err)
		return
	}

	if err := hltClient.PutMessage(msg); err != nil {
		t.Error(err)
		return
	}

	hashes, err := hltClient.GetHashes()
	if err != nil {
		t.Error(err)
		return
	}

	if len(hashes) != 1 {
		t.Error("len hashes != 1")
		return
	}

	if hashes[0] != encoding.HexEncode(msg.GetBody().GetHash()) {
		t.Error("hashes not equals")
		return
	}
}
