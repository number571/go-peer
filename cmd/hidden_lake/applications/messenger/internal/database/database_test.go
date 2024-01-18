package database

import (
	"bytes"
	"os"
	"testing"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/pkg/settings"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/storage"
	testutils "github.com/number571/go-peer/test/utils"
)

const (
	tcPath = "database.db"
)

func TestDatabase(t *testing.T) {
	t.Parallel()

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
	pseudonym := random.NewStdPRNG().GetString(settings.CPseudonymSize)
	err1 := db.Push(rel, NewMessage(true, pseudonym, []byte(testutils.TcBody)))
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

	if msgs[0].GetPseudonym() != pseudonym {
		t.Error("msgs[0].GetSenderID() != pseudonym")
		return
	}

	if !bytes.Equal(msgs[0].GetMessage(), []byte(testutils.TcBody)) {
		t.Error("!bytes.Equal(msgs[0].GetMessage(), []byte(testutils.TcBody))")
		return
	}
}
