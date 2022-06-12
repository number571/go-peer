package database

import (
	"bytes"
	"os"
	"testing"

	"github.com/number571/go-peer/crypto/hashing"
	"github.com/number571/go-peer/offline/message"
)

const (
	tcPathDB           = "test.db"
	tcMessageTitle     = "test-title"
	tcMessageBody      = "test-body"
	tcMessageRawConcat = tcMessageTitle + tcMessageBody
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
		t.Errorf("init size != 0")
		return
	}

	msg := message.NewMessage([]byte(tcMessageTitle), []byte(tcMessageBody))
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
