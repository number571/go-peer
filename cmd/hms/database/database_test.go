package database

import (
	"bytes"
	"os"
	"testing"

	"github.com/number571/go-peer/client"
	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/crypto/hashing"
	"github.com/number571/go-peer/payload"
	"github.com/number571/go-peer/testutils"
)

const (
	tcPathDB = "test.db"
)

var (
	tgKey = hashing.NewSHA256Hasher([]byte("test-key")).Bytes()
	tgDB  IKeyValueDB
)

func testHmsDefaultInit(dbPath string) {
	os.RemoveAll(dbPath)
	tgDB = NewKeyValueDB(dbPath)
}

func TestDB(t *testing.T) {
	testHmsDefaultInit(tcPathDB)
	defer func() {
		tgDB.Close()
		os.RemoveAll(tcPathDB)
	}()

	size, err := tgDB.Size(tgKey)
	if err != nil {
		t.Error(err)
		return
	}
	if size != 0 {
		t.Error("init size != 0")
		return
	}

	cl := client.NewClient(
		client.NewSettings(10, (1<<20)),
		asymmetric.NewRSAPrivKey(1024),
	)

	msg, err := cl.Encrypt(
		cl.PubKey(),
		payload.NewPayload(uint64(testutils.TcHead), []byte(testutils.TcBody)),
	)
	if err != nil {
		t.Error(err)
		return
	}

	err = tgDB.Push(tgKey, msg)
	if err != nil {
		t.Error(err)
		return
	}

	size, err = tgDB.Size(tgKey)
	if err != nil {
		t.Error(err)
		return
	}
	if size != 1 {
		t.Error("after push size != 1")
		return
	}

	loadMsg, err := tgDB.Load(tgKey, 0)
	if err != nil {
		t.Error(err)
		return
	}

	pubKey, pl, err := cl.Decrypt(loadMsg)
	if err != nil {
		t.Error(err)
		return
	}

	if pubKey.Address().String() != cl.PubKey().Address().String() {
		t.Error("load public key != init public key")
		return
	}

	if pl.Head() != uint64(testutils.TcHead) {
		t.Error("load msg head != init head")
		return
	}

	if !bytes.Equal(pl.Body(), []byte(testutils.TcBody)) {
		t.Error("load msg body != init body")
		return
	}

	err = tgDB.Clean()
	if err != nil {
		t.Error(err)
		return
	}

	size, err = tgDB.Size(tgKey)
	if err != nil {
		t.Error(err)
		return
	}
	if size != 0 {
		t.Error("after clean size != 0")
		return
	}

	err = tgDB.Close()
	if err != nil {
		t.Error(err)
		return
	}
}
