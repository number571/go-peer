package handler

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	testutils "github.com/number571/go-peer/test/_data"
)

func TestHandleHashesAPI(t *testing.T) {
	t.Parallel()

	addr := testutils.TgAddrs[19]
	os.RemoveAll(fmt.Sprintf(databaseTemplate, addr))

	srv, connKeeper, cancel, wDB, hltClient := testAllRun(addr, "")
	defer testAllFree(addr, srv, connKeeper, cancel, wDB)

	privKey := asymmetric.LoadRSAPrivKey(testutils.Tc1PrivKey1024)
	pubKey := asymmetric.LoadRSAPubKey(testutils.TgPubKeys[0])

	client := client.NewClient(
		message.NewSettings(&message.SSettings{
			FMessageSizeBytes: testutils.TCMessageSize,
			FKeySizeBits:      testutils.TcKeySize,
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

	netMsg := net_message.NewMessage(
		testNetworkMessageSettings(),
		payload.NewPayload(hls_settings.CNetworkMask, msg.ToBytes()),
	)
	if err := hltClient.PutMessage(netMsg); err != nil {
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

	if !bytes.Equal(encoding.HexDecode(hashes[0]), netMsg.GetHash()) {
		t.Error("hashes not equals")
		return
	}
}
