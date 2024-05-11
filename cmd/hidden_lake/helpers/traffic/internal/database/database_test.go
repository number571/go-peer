package database

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	testutils "github.com/number571/go-peer/test/utils"
)

const (
	tcPathDBTemplate = "test_database_%d.db"
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

func TestSettings(t *testing.T) {
	t.Parallel()

	for i := 0; i < 2; i++ {
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
			FMessagesCapacity: testutils.TCCapacity,
		})
	case 1:
		_ = NewSettings(&SSettings{
			FPath: dbPath,
		})
	}
}

// // The test fails when the user is root.
// func TestIInitDatabase(t *testing.T) {
// 	t.Parallel()

// 	prng := random.NewStdPRNG()
// 	path := "/" + prng.GetString(32) + "/" + prng.GetString(32) + "/" + prng.GetString(32)
// 	defer os.RemoveAll(path)

// 	_, err := NewDatabase(NewSettings(&SSettings{
// 		FPath:             path,
// 		FNetworkKey:       testutils.TCNetworkKey,
// 		FWorkSizeBits:     testutils.TCWorkSize,
// 		FMessagesCapacity: testutils.TCCapacity,
// 	}))
// 	if err == nil {
// 		t.Error("success init database with invalid path")
// 		return
// 	}
// }

func TestDatabaseLoad(t *testing.T) {
	t.Parallel()

	testDatabaseLoad(t, 8, NewDatabase)
	testDatabaseLoad(t, 9, NewInMemoryDatabase)
}

func testDatabaseLoad(t *testing.T, numDB int, dbConstruct func(pSett ISettings) (IDatabase, error)) {
	pathDB := fmt.Sprintf(tcPathDBTemplate, numDB)
	os.RemoveAll(pathDB)

	kvDB, err := dbConstruct(NewSettings(&SSettings{
		FPath:             pathDB,
		FNetworkKey:       testutils.TCNetworkKey,
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

	if _, err := kvDB.Load([]byte("abc")); err == nil {
		t.Error("success load not exist message (incorrect)")
		return
	}

	hash := hashing.NewSHA256Hasher([]byte{123}).ToBytes()
	_, errLoad := kvDB.Load(hash)
	if errLoad == nil {
		t.Error("success load not exist message (hash)")
		return
	}

	if !errors.Is(errLoad, ErrMessageIsNotExist) {
		t.Error("got incorrect error type (load)")
		return
	}

	ptrDB, ok := kvDB.(*sDatabase)
	if !ok {
		// pass inmemory_database
		return
	}

	hash2 := hashing.NewSHA256Hasher([]byte{123}).ToBytes()
	if err := ptrDB.fDB.Set(getKeyMessage(hash2), []byte{123}); err != nil {
		t.Error(err)
		return
	}

	if _, err := kvDB.Load(hash2); err == nil {
		t.Error("success load with incorrect hash")
		return
	}
}

func TestDatabaseHashes(t *testing.T) {
	t.Parallel()

	testDatabaseHashes(t, 4, NewDatabase)
	testDatabaseHashes(t, 5, NewInMemoryDatabase)
}

func testDatabaseHashes(t *testing.T, numDB int, dbConstruct func(pSett ISettings) (IDatabase, error)) {
	const (
		messagesCapacity = 3
	)

	pathDB := fmt.Sprintf(tcPathDBTemplate, numDB)
	os.RemoveAll(pathDB)

	kvDB, err := dbConstruct(NewSettings(&SSettings{
		FPath:             pathDB,
		FNetworkKey:       testutils.TCNetworkKey,
		FWorkSizeBits:     testutils.TCWorkSize,
		FMessagesCapacity: messagesCapacity,
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
			FMessageSizeBytes: testutils.TCMessageSize,
			FKeySizeBits:      testutils.TcKeySize,
		}),
		asymmetric.LoadRSAPrivKey(testutils.Tc1PrivKey1024),
	)

	pushHashes := make([][]byte, 0, messagesCapacity+1)
	for i := 0; i < messagesCapacity+1; i++ {
		msg, err := newNetworkMessageWithData(cl, testutils.TCNetworkKey, strconv.Itoa(i))
		if err != nil {
			t.Error(err)
			return
		}
		if err := kvDB.Push(msg); err != nil {
			t.Error(err)
			return
		}
		if kvDB.Pointer() != uint64(i+1)%messagesCapacity {
			t.Error("got invalid pointer")
			return
		}
		pushHashes = append(pushHashes, msg.GetHash())
	}

	for i := uint64(0); i < messagesCapacity+1; i++ {
		hash, err := kvDB.Hash(i)
		if err != nil {
			break
		}
		if bytes.Equal(hash, pushHashes[0]) {
			t.Error("hash not overwritten")
			return
		}
		index := (2 + i) % (messagesCapacity)
		if !bytes.Equal(hash, pushHashes[1:][index]) {
			t.Error("got invalid hash")
			return
		}
	}
}

func TestDatabasePush(t *testing.T) {
	t.Parallel()

	testDatabasePush(t, 6, NewDatabase)
	testDatabasePush(t, 7, NewInMemoryDatabase)
}

func testDatabasePush(t *testing.T, numDB int, dbConstruct func(pSett ISettings) (IDatabase, error)) {
	pathDB := fmt.Sprintf(tcPathDBTemplate, numDB)
	os.RemoveAll(pathDB)

	kvDB, err := dbConstruct(NewSettings(&SSettings{
		FPath:             pathDB,
		FNetworkKey:       testutils.TCNetworkKey,
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
			FMessageSizeBytes: (10 << 10),
			FKeySizeBits:      testutils.TcKeySize,
		}),
		asymmetric.LoadRSAPrivKey(testutils.Tc1PrivKey1024),
	)

	msgTest, err := newNetworkMessage(clTest, "some-another-key")
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
			FMessageSizeBytes: testutils.TCMessageSize,
			FKeySizeBits:      testutils.TcKeySize,
		}),
		asymmetric.LoadRSAPrivKey(testutils.Tc1PrivKey1024),
	)

	msg1, err := newNetworkMessage(cl, testutils.TCNetworkKey)
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

	if !errors.Is(errPush, ErrMessageIsExist) {
		t.Error("got incorrect error type (push)")
		return
	}

	msg2, err := newNetworkMessage(cl, testutils.TCNetworkKey)
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

	kvDB, err := NewDatabase(NewSettings(&SSettings{
		FPath:             pathDB,
		FNetworkKey:       testutils.TCNetworkKey,
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
			FMessageSizeBytes: testutils.TCMessageSize,
			FKeySizeBits:      testutils.TcKeySize,
		}),
		asymmetric.LoadRSAPrivKey(testutils.Tc1PrivKey1024),
	)

	putHashes := make([][]byte, 0, 3)
	for i := 0; i < 3; i++ {
		msg, err := newNetworkMessage(cl, testutils.TCNetworkKey)
		if err != nil {
			t.Error(err)
			return
		}

		if err := kvDB.Push(msg); err != nil {
			t.Error(err)
			return
		}
		putHashes = append(putHashes, msg.GetHash())
	}

	getHashes := make([][]byte, 0, 3)
	for i := uint64(0); ; i++ {
		hash, err := kvDB.Hash(i)
		if err != nil {
			break
		}
		getHashes = append(getHashes, hash)
	}

	if len(getHashes) != 3 {
		t.Error("len getHashes != 3")
		return
	}

	for i := range getHashes {
		if !bytes.Equal(getHashes[i], putHashes[i]) {
			t.Errorf("getHashes[%d] != putHashes[%d]", i, i)
			return
		}
	}

	for _, getHash := range getHashes {
		loadNetMsg, err := kvDB.Load(getHash)
		if err != nil {
			t.Error(err)
			return
		}

		msgHash := loadNetMsg.GetHash()
		if !bytes.Equal(getHash, msgHash) {
			t.Errorf("getHash[%s] != msgHash[%s]", getHash, msgHash)
			return
		}

		loadMsg, err := message.LoadMessage(cl.GetSettings(), loadNetMsg.GetPayload().GetBody())
		if err != nil {
			t.Error(err)
			return
		}

		pubKey, pl, err := cl.DecryptMessage(loadMsg)
		if err != nil {
			t.Error(err)
			return
		}

		if pubKey.GetHasher().ToString() != cl.GetPubKey().GetHasher().ToString() {
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

func newNetworkMessageWithData(cl client.IClient, networkKey, data string) (net_message.IMessage, error) {
	msg, err := cl.EncryptPayload(
		cl.GetPubKey(),
		payload.NewPayload64(uint64(testutils.TcHead), []byte(data)),
	)
	if err != nil {
		return nil, err
	}
	netMsg := net_message.NewMessage(
		net_message.NewSettings(&net_message.SSettings{
			FNetworkKey:   networkKey,
			FWorkSizeBits: testutils.TCWorkSize,
		}),
		payload.NewPayload64(0, msg.ToBytes()),
	)
	return netMsg, nil
}

func newNetworkMessage(cl client.IClient, networkKey string) (net_message.IMessage, error) {
	return newNetworkMessageWithData(cl, networkKey, testutils.TcBody)
}
