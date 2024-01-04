package database

import (
	"bytes"
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/storage"
	testutils "github.com/number571/go-peer/test/_data"
)

const (
	tcPath = "database.db"
)

func TestDatabase(t *testing.T) {
	os.RemoveAll(tcPath)
	defer os.RemoveAll(tcPath)

	db, err := NewKeyValueDB(storage.NewSettings(&storage.SSettings{
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
	err1 := db.Push(rel, NewMessage(true, friend.GetHasher().ToBytes(), []byte(testutils.TcBody)))
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

	if msgs[0].GetSenderID() != friend.GetHasher().ToString() {
		t.Error("msgs[0].GetSenderID() != friend.GetHasher().ToString()")
		return
	}

	if !bytes.Equal(msgs[0].GetMessage(), []byte(testutils.TcBody)) {
		t.Error("!bytes.Equal(msgs[0].GetMessage(), []byte(testutils.TcBody))")
		return
	}
}
