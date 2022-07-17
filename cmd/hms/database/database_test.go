package database

import (
	"bytes"
	"os"
	"testing"

	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/crypto/hashing"
	"github.com/number571/go-peer/local/client"
	"github.com/number571/go-peer/local/payload"
	"github.com/number571/go-peer/local/routing"
	"github.com/number571/go-peer/settings"
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
		settings.NewSettings(),
		asymmetric.NewRSAPrivKey(1024),
	)

	err = tgDB.Push(tgKey, cl.Encrypt(
		routing.NewRoute(cl.PubKey()),
		payload.NewPayload(uint64(testutils.TcHead), []byte(testutils.TcBody)),
	))
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

	pubKey, pl := cl.Decrypt(loadMsg)
	if pubKey == nil {
		t.Error("load message is nil")
		return
	}

	if pl.Head() != uint64(testutils.TcHead) {
		t.Error("load msg head != init head")
	}

	if !bytes.Equal(pl.Body(), []byte(testutils.TcBody)) {
		t.Error("load msg body != init body")
	}

	err = tgDB.Clean()
	if err != nil {
		t.Error(err)
	}

	size, err = tgDB.Size(tgKey)
	if err != nil {
		t.Error(err)
		return
	}
	if size != 0 {
		t.Error("after clean size != 0")
	}

	err = tgDB.Close()
	if err != nil {
		t.Error(err)
	}
}
