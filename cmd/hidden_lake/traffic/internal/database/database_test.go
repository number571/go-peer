package database

import (
	"bytes"
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload"
	testutils "github.com/number571/go-peer/test/_data"
)

const (
	tcPathDB = "test.db"
)

var (
	tgDB IKVDatabase
)

func testHmsDefaultInit(dbPath string) error {
	os.RemoveAll(dbPath)
	var err error
	tgDB, err = NewKeyValueDB(NewSettings(&SSettings{
		FPath:             dbPath,
		FMessageSizeBytes: testutils.TCMessageSize,
		FWorkSizeBits:     testutils.TCWorkSize,
		FCapacity:         testutils.TCCapacity,
	}))
	if err != nil {
		return err
	}
	return nil
}

func TestDB(t *testing.T) {
	err := testHmsDefaultInit(tcPathDB)
	if err != nil {
		t.Error(err)
		return
	}

	defer func() {
		tgDB.Close()
		os.RemoveAll(tcPathDB)
	}()

	cl := client.NewClient(
		message.NewSettings(&message.SSettings{
			FWorkSizeBits:     testutils.TCWorkSize,
			FMessageSizeBytes: testutils.TCMessageSize,
		}),
		asymmetric.LoadRSAPrivKey(testutils.TcPrivKey),
	)

	putHashes := make([]string, 0, 3)

	for i := 0; i < 3; i++ {
		msg, err := cl.EncryptPayload(
			cl.GetPubKey(),
			payload.NewPayload(uint64(testutils.TcHead), []byte(testutils.TcBody)),
		)
		if err != nil {
			t.Error(err)
			return
		}

		if err := tgDB.Push(msg); err != nil {
			t.Error(err)
			return
		}
		putHashes = append(putHashes, encoding.HexEncode(msg.GetBody().GetHash()))
	}

	getHashes, err := tgDB.Hashes()
	if err != nil {
		t.Error(err)
		return
	}

	if len(getHashes) != 3 {
		t.Error("len getHashes != 3")
		return
	}

	for i := range getHashes {
		if getHashes[i] != putHashes[i] {
			t.Errorf("getHashes[%d] != putHashes[%d]", i, i)
			return
		}
	}

	for _, getHash := range getHashes {
		loadMsg, err := tgDB.Load(getHash)
		if err != nil {
			t.Error(err)
			return
		}

		msgHash := encoding.HexEncode(loadMsg.GetBody().GetHash())
		if getHash != msgHash {
			t.Errorf("getHash[%s] != msgHash[%s]", getHash, msgHash)
			return
		}

		pubKey, pl, err := cl.DecryptMessage(loadMsg)
		if err != nil {
			t.Error(err)
			return
		}

		if pubKey.GetAddress().ToString() != cl.GetPubKey().GetAddress().ToString() {
			t.Error("load public key != init public key")
			return
		}

		if pl.GetHead() != uint64(testutils.TcHead) {
			t.Error("load msg head != init head")
			return
		}

		if !bytes.Equal(pl.GetBody(), []byte(testutils.TcBody)) {
			t.Error("load msg body != init body")
			return
		}
	}

	if err := tgDB.Close(); err != nil {
		t.Error(err)
		return
	}
}
