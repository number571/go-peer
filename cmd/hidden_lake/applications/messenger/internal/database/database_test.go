package database

import (
	"bytes"
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/database"
	testutils "github.com/number571/go-peer/test/utils"
)

const (
	tcPath = "database.db"
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SDatabaseError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestDatabase(t *testing.T) {
	t.Parallel()

	os.RemoveAll(tcPath)
	defer os.RemoveAll(tcPath)

	db, err := NewKeyValueDB(database.NewSettings(&database.SSettings{
		FPath: tcPath,
	}))
	if err != nil {
		t.Error(err)
		return
	}
	defer db.Close()

	iam := asymmetric.LoadRSAPubKey(testutils.TgPubKeys[0])
	friend := asymmetric.LoadRSAPubKey(testutils.TgPubKeys[1])

	rel := NewRelation(iam, friend)
	err1 := db.Push(rel, NewMessage(true, []byte(testutils.TcBody)))
	if err1 != nil {
		t.Error(err1)
		return
	}

	size := db.Size(rel)
	if size != 1 {
		t.Error("size != 1")
		return
	}

	msgs, err := db.Load(rel, 0, size)
	if err != nil {
		t.Error(err)
		return
	}

	if len(msgs) != 1 {
		t.Error("len(msgs) != 1")
		return
	}

	if !msgs[0].IsIncoming() {
		t.Error("!msgs[0].IsIncoming()")
		return
	}

	if !bytes.Equal(msgs[0].GetMessage(), []byte(testutils.TcBody)) {
		t.Error("!bytes.Equal(msgs[0].GetMessage(), []byte(testutils.TcBody))")
		return
	}
}
