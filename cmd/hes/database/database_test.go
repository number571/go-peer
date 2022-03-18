package database

import (
	"bytes"
	"os"
	"testing"

	"github.com/number571/go-peer/local"
)

const (
	tcPathDB = "test_hes.db"
	tcKey    = "test-key"

	tcMessageTitle     = "test-title"
	tcMessageBody      = "test-body"
	tcMessageRawConcat = tcMessageTitle + tcMessageBody
)

var (
	tgDB IKeyValueDB
)

func testHesDefaultInit(path string) {
	tgDB = NewKeyValueDB(path)
}

func TestDB(t *testing.T) {
	testHesDefaultInit(tcPathDB)
	defer os.RemoveAll(tcPathDB)

	if tgDB.Size([]byte(tcKey)) != 0 {
		t.Errorf("init size != 0")
	}

	msg := local.NewMessage([]byte(tcMessageTitle), []byte(tcMessageBody))
	err := tgDB.Push([]byte(tcKey), msg)
	if err != nil {
		t.Error(err)
	}

	if tgDB.Size([]byte(tcKey)) != 1 {
		t.Errorf("after push size != 1")
	}

	loadMsg := tgDB.Load([]byte(tcKey), 0)
	if loadMsg == nil {
		t.Errorf("load msg is nil")
	}

	if !bytes.Equal(loadMsg.Body().Data(), []byte(msg.Body().Data())) {
		t.Errorf("load msg (title||body) != init (title||body)")
	}

	err = tgDB.Clean()
	if err != nil {
		t.Error(err)
	}

	if tgDB.Size([]byte(tcKey)) != 0 {
		t.Errorf("after clean size != 0")
	}

	err = tgDB.Close()
	if err != nil {
		t.Error(err)
	}
}
