package handler

import (
	"fmt"
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload"
	testutils "github.com/number571/go-peer/test/_data"
)

func TestHandleHashesAPI(t *testing.T) {
	addr := testutils.TgAddrs[19]
	os.RemoveAll(fmt.Sprintf(databaseTemplate, addr))

	srv, connKeeper, wDB, hltClient := testAllRun(addr, "")
	defer testAllFree(addr, srv, connKeeper, wDB)

	privKey := asymmetric.LoadRSAPrivKey(testutils.TcPrivKey)
	pubKey := asymmetric.LoadRSAPubKey(testutils.TgPubKeys[0])

	client := client.NewClient(
		message.NewSettings(&message.SSettings{
			FMessageSizeBytes: testutils.TCMessageSize,
			FWorkSizeBits:     testutils.TCWorkSize,
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
