package database

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/payload"
	testutils "github.com/number571/go-peer/test/_data"
)

const (
	tcPathDBTemplate = "test_database_%d.db"
)

func TestSettings(t *testing.T) {
	t.Parallel()

	for i := 0; i < 3; i++ {
		testSettings(t, i)
	}
}

func testSettings(t *testing.T, n int) {
	dbPath := "test_settings.db"
	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()

	switch n {
	case 0:
		_ = NewSettings(&SSettings{
			FMessageSizeBytes: testutils.TCMessageSize,
			FWorkSizeBits:     testutils.TCWorkSize,
			FMessagesCapacity: testutils.TCCapacity,
		})
	case 1:
		_ = NewSettings(&SSettings{
			FPath:             dbPath,
			FWorkSizeBits:     testutils.TCWorkSize,
			FMessagesCapacity: testutils.TCCapacity,
		})
	case 2:
		_ = NewSettings(&SSettings{
			FPath:             dbPath,
			FMessageSizeBytes: testutils.TCMessageSize,
			FWorkSizeBits:     testutils.TCWorkSize,
		})
	}
}

func TestIInitDatabase(t *testing.T) {
	t.Parallel()

	prng := random.NewStdPRNG()
	path := "/" + prng.GetString(32) + "/" + prng.GetString(32) + "/" + prng.GetString(32)

	_, err := NewKeyValueDB(NewSettings(&SSettings{
		FPath:             path,
		FMessageSizeBytes: testutils.TCMessageSize,
		FWorkSizeBits:     testutils.TCWorkSize,
		FMessagesCapacity: testutils.TCCapacity,
	}))
	if err == nil {
		t.Error("success init database with invalid path")
		return
	}
}

func TestDatabaseLoadPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()

	t.Parallel()

	pathDB := fmt.Sprintf(tcPathDBTemplate, 4)
	os.RemoveAll(pathDB)

	kvDB, err := NewKeyValueDB(NewSettings(&SSettings{
		FPath:             pathDB,
		FMessageSizeBytes: testutils.TCMessageSize,
		FWorkSizeBits:     testutils.TCWorkSize,
		FMessagesCapacity: testutils.TCCapacity,
	}))
	if err != nil {
		t.Error(err)
		return
	}

	defer func() {
		kvDB.Close()
		os.RemoveAll(pathDB)
	}()

	ptrDB := kvDB.(*sKeyValueDB)

	hash := hashing.NewSHA256Hasher([]byte{123}).ToBytes()
	if err := ptrDB.fDB.Set(getKeyMessage(hash), []byte{123}); err != nil {
		t.Error(err)
		return
	}

	_, _ = kvDB.Load(encoding.HexEncode(hash)) // panic
}

func TestDatabaseHashesPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()

	t.Parallel()

	pathDB := fmt.Sprintf(tcPathDBTemplate, 3)
	os.RemoveAll(pathDB)

	kvDB, err := NewKeyValueDB(NewSettings(&SSettings{
		FPath:             pathDB,
		FMessageSizeBytes: testutils.TCMessageSize,
		FWorkSizeBits:     testutils.TCWorkSize,
		FMessagesCapacity: testutils.TCCapacity,
	}))
	if err != nil {
		t.Error(err)
		return
	}

	defer func() {
		kvDB.Close()
		os.RemoveAll(pathDB)
	}()

	ptrDB := kvDB.(*sKeyValueDB)

	if err := ptrDB.fDB.Set(getKeyHash(0), []byte{123}); err != nil {
		t.Error(err)
		return
	}

	_, _ = kvDB.Hashes() // panic
}

func TestDatabaseLoad(t *testing.T) {
	t.Parallel()

	pathDB := fmt.Sprintf(tcPathDBTemplate, 2)
	os.RemoveAll(pathDB)

	kvDB, err := NewKeyValueDB(NewSettings(&SSettings{
		FPath:             pathDB,
		FMessageSizeBytes: testutils.TCMessageSize,
		FWorkSizeBits:     testutils.TCWorkSize,
		FMessagesCapacity: testutils.TCCapacity,
	}))
	if err != nil {
		t.Error(err)
		return
	}

	defer func() {
		kvDB.Close()
		os.RemoveAll(pathDB)
	}()

	if _, err := kvDB.Load("abc"); err == nil {
		t.Error("success load not exist message (incorrect)")
		return
	}

	hash := hashing.NewSHA256Hasher([]byte{123}).ToBytes()
	hashStr := encoding.HexEncode(hash)

	_, errLoad := kvDB.Load(hashStr)
	if errLoad == nil {
		t.Error("success load not exist message (hash)")
		return
	}

	if !errors.HasError(errLoad, &SIsNotExistError{}) {
		t.Error("got incorrect error type (load)")
		return
	}
}

func TestDatabasePush(t *testing.T) {
	t.Parallel()

	pathDB := fmt.Sprintf(tcPathDBTemplate, 1)
	os.RemoveAll(pathDB)

	kvDB, err := NewKeyValueDB(NewSettings(&SSettings{
		FPath:             pathDB,
		FMessageSizeBytes: testutils.TCMessageSize,
		FWorkSizeBits:     testutils.TCWorkSize,
		FMessagesCapacity: 1,
	}))
	if err != nil {
		t.Error(err)
		return
	}

	defer func() {
		kvDB.Close()
		os.RemoveAll(pathDB)
	}()

	clTest := client.NewClient(
		message.NewSettings(&message.SSettings{
			FWorkSizeBits:     1,
			FMessageSizeBytes: (10 << 10),
		}),
		asymmetric.LoadRSAPrivKey(testutils.Tc1PrivKey1024),
	)

	msgTest, err := newMessage(clTest)
	if err != nil {
		t.Error(err)
		return
	}

	if err := kvDB.Push(msgTest); err == nil {
		t.Error("success push message with difference setting")
		return
	}

	cl := client.NewClient(
		message.NewSettings(&message.SSettings{
			FWorkSizeBits:     testutils.TCWorkSize,
			FMessageSizeBytes: testutils.TCMessageSize,
		}),
		asymmetric.LoadRSAPrivKey(testutils.Tc1PrivKey1024),
	)

	msg1, err := newMessage(cl)
	if err != nil {
		t.Error(err)
		return
	}

	if err := kvDB.Push(msg1); err != nil {
		t.Error(err)
		return
	}

	errPush := kvDB.Push(msg1)
	if errPush == nil {
		t.Error("success push duplicate")
		return
	}

	if !errors.HasError(errPush, &SIsExistError{}) {
		t.Error("got incorrect error type (push)")
		return
	}

	msg2, err := newMessage(cl)
	if err != nil {
		t.Error(err)
		return
	}

	if err := kvDB.Push(msg2); err != nil {
		t.Error(err)
		return
	}
}

func TestDatabase(t *testing.T) {
	t.Parallel()

	pathDB := fmt.Sprintf(tcPathDBTemplate, 0)
	os.RemoveAll(pathDB)

	kvDB, err := NewKeyValueDB(NewSettings(&SSettings{
		FPath:             pathDB,
		FMessageSizeBytes: testutils.TCMessageSize,
		FWorkSizeBits:     testutils.TCWorkSize,
		FMessagesCapacity: testutils.TCCapacity,
	}))
	if err != nil {
		t.Error(err)
		return
	}

	defer func() {
		kvDB.Close()
		os.RemoveAll(pathDB)
	}()

	cl := client.NewClient(
		message.NewSettings(&message.SSettings{
			FWorkSizeBits:     testutils.TCWorkSize,
			FMessageSizeBytes: testutils.TCMessageSize,
		}),
		asymmetric.LoadRSAPrivKey(testutils.Tc1PrivKey1024),
	)

	putHashes := make([]string, 0, 3)
	for i := 0; i < 3; i++ {
		msg, err := newMessage(cl)
		if err != nil {
			t.Error(err)
			return
		}

		if err := kvDB.Push(msg); err != nil {
			t.Error(err)
			return
		}
		putHashes = append(putHashes, encoding.HexEncode(msg.GetBody().GetHash()))
	}

	getHashes, err := kvDB.Hashes()
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
		loadMsg, err := kvDB.Load(getHash)
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

	if err := kvDB.Close(); err != nil {
		t.Error(err)
		return
	}
}

func newMessage(cl client.IClient) (message.IMessage, error) {
	msg, err := cl.EncryptPayload(
		cl.GetPubKey(),
		payload.NewPayload(uint64(testutils.TcHead), []byte(testutils.TcBody)),
	)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
