package database

import (
	"bytes"
	"os"
	"testing"

	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/local"
)

const (
	tcPathDB = "database_test.db"
)

const (
	tcMessageTitle     = "test-title"
	tcMessageBody      = "test-body"
	tcMessageRawConcat = tcMessageTitle + tcMessageBody
)

var (
	tgKey = crypto.NewHasher([]byte("test-key")).Bytes()
	tgDB  IKeyValueDB
)

func testHmsDefaultInit(path string) {
	tgDB = NewKeyValueDB(path)
}

func TestDB(t *testing.T) {
	testHmsDefaultInit(tcPathDB)
	defer os.RemoveAll(tcPathDB)

	size, err := tgDB.Size(tgKey)
	if err != nil {
		t.Error(err)
		return
	}
	if size != 0 {
		t.Errorf("init size != 0")
		return
	}

	msg := local.NewMessage([]byte(tcMessageTitle), []byte(tcMessageBody))
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
		t.Errorf("after push size != 1")
		return
	}

	loadMsg, err := tgDB.Load(tgKey, 0)
	if err != nil {
		t.Error(err)
		return
	}

	if !bytes.Equal(loadMsg.Body().Data(), []byte(msg.Body().Data())) {
		t.Errorf("load msg (title||body) != init (title||body)")
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
		t.Errorf("after clean size != 0")
	}

	err = tgDB.Close()
	if err != nil {
		t.Error(err)
	}
}
