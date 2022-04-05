package database

import (
	"os"
	"testing"
)

const (
	tcPathDB = "database_test.db"
)

const (
	tcHash             = "test-hash"
	tcMessageTitle     = "test-title"
	tcMessageBody      = "test-body"
	tcMessageRawConcat = tcMessageTitle + tcMessageBody
)

var (
	tgDB IKeyValueDB
)

func testHmsDefaultInit(path string) {
	tgDB = NewKeyValueDB(path)
}

func TestDB(t *testing.T) {
	testHmsDefaultInit(tcPathDB)
	defer os.RemoveAll(tcPathDB)

	err := tgDB.Push([]byte(tcHash))
	if err != nil {
		t.Error(err)
	}

	exist := tgDB.Exist([]byte(tcHash))
	if !exist {
		t.Errorf("load msg is nil")
	}

	err = tgDB.Clean()
	if err != nil {
		t.Error(err)
	}

	err = tgDB.Close()
	if err != nil {
		t.Error(err)
	}
}
