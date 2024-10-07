package storage

import (
	"bytes"
	"errors"
	"strconv"
	"testing"

	hlt_database "github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/internal/database"
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/storage/cache"
	testutils "github.com/number571/go-peer/test/utils"
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SStorageError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestStorageLoad(t *testing.T) {
	t.Parallel()

	db := NewMessageStorage(
		net_message.NewSettings(&net_message.SSettings{
			FWorkSizeBits: testutils.TCWorkSize,
			FNetworkKey:   testutils.TCNetworkKey,
		}),
		hlt_database.NewVoidKVDatabase(),
		cache.NewLRUCache(testutils.TCCapacity),
	)

	if _, err := db.Load([]byte("abc")); err == nil {
		t.Error("success load not exist message (incorrect)")
		return
	}

	hash := hashing.NewSHA256Hasher([]byte{123}).ToBytes()
	_, errLoad := db.Load(hash)
	if errLoad == nil {
		t.Error("success load not exist message (hash)")
		return
	}

	if !errors.Is(errLoad, ErrMessageIsNotExist) {
		t.Error("got incorrect error type (load)")
		return
	}
}

func TestStorageHashes(t *testing.T) {
	t.Parallel()
	const messagesCapacity = 3

	db := NewMessageStorage(
		net_message.NewSettings(&net_message.SSettings{
			FWorkSizeBits: testutils.TCWorkSize,
			FNetworkKey:   testutils.TCNetworkKey,
		}),
		hlt_database.NewVoidKVDatabase(),
		cache.NewLRUCache(messagesCapacity),
	)

	cl := client.NewClient(
		message.NewSettings(&message.SSettings{
			FMessageSizeBytes: testutils.TCMessageSize,
			FKeySizeBits:      testutils.TcKeySize,
		}),
		asymmetric.LoadRSAPrivKey(testutils.TcPrivKey1024),
	)

	pushHashes := make([][]byte, 0, messagesCapacity+1)
	for i := 0; i < messagesCapacity+1; i++ {
		msg, err := newNetworkMessageWithData(cl, testutils.TCNetworkKey, strconv.Itoa(i))
		if err != nil {
			t.Error(err)
			return
		}
		if err := db.Push(msg); err != nil {
			t.Error(err)
			return
		}
		if db.Pointer() != uint64(i+1)%messagesCapacity {
			t.Error("got invalid pointer")
			return
		}
		pushHashes = append(pushHashes, msg.GetHash())
	}

	for i := uint64(0); i < messagesCapacity+1; i++ {
		hash, err := db.Hash(i)
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

func TestStoragePush(t *testing.T) {
	t.Parallel()

	db := NewMessageStorage(
		net_message.NewSettings(&net_message.SSettings{
			FWorkSizeBits: testutils.TCWorkSize,
			FNetworkKey:   testutils.TCNetworkKey,
		}),
		hlt_database.NewVoidKVDatabase(),
		cache.NewLRUCache(1),
	)

	clTest := client.NewClient(
		message.NewSettings(&message.SSettings{
			FMessageSizeBytes: (10 << 10),
			FKeySizeBits:      testutils.TcKeySize,
		}),
		asymmetric.LoadRSAPrivKey(testutils.TcPrivKey1024),
	)

	msgTest, err := newNetworkMessage(clTest, "some-another-key")
	if err != nil {
		t.Error(err)
		return
	}

	if err := db.Push(msgTest); err == nil {
		t.Error("success push message with difference setting")
		return
	}

	cl := client.NewClient(
		message.NewSettings(&message.SSettings{
			FMessageSizeBytes: testutils.TCMessageSize,
			FKeySizeBits:      testutils.TcKeySize,
		}),
		asymmetric.LoadRSAPrivKey(testutils.TcPrivKey1024),
	)

	msg1, err := newNetworkMessage(cl, testutils.TCNetworkKey)
	if err != nil {
		t.Error(err)
		return
	}

	if err := db.Push(msg1); err != nil {
		t.Error(err)
		return
	}

	errPush := db.Push(msg1)
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

	if err := db.Push(msg2); err != nil {
		t.Error(err)
		return
	}
}

func TestStorage(t *testing.T) {
	t.Parallel()

	db := NewMessageStorage(
		net_message.NewSettings(&net_message.SSettings{
			FWorkSizeBits: testutils.TCWorkSize,
			FNetworkKey:   testutils.TCNetworkKey,
		}),
		hlt_database.NewVoidKVDatabase(),
		cache.NewLRUCache(4),
	)

	cl := client.NewClient(
		message.NewSettings(&message.SSettings{
			FMessageSizeBytes: testutils.TCMessageSize,
			FKeySizeBits:      testutils.TcKeySize,
		}),
		asymmetric.LoadRSAPrivKey(testutils.TcPrivKey1024),
	)

	putHashes := make([][]byte, 0, 3)
	for i := 0; i < 3; i++ {
		msg, err := newNetworkMessage(cl, testutils.TCNetworkKey)
		if err != nil {
			t.Error(err)
			return
		}
		if err := db.Push(msg); err != nil {
			t.Error(err)
			return
		}
		putHashes = append(putHashes, msg.GetHash())
	}

	getHashes := make([][]byte, 0, 3)
	for i := uint64(0); ; i++ {
		hash, err := db.Hash(i)
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
		loadNetMsg, err := db.Load(getHash)
		if err != nil {
			t.Error(err)
			return
		}

		msgHash := loadNetMsg.GetHash()
		if !bytes.Equal(getHash, msgHash) {
			t.Errorf("getHash[%s] != msgHash[%s]", getHash, msgHash)
			return
		}

		pubKey, decMsg, err := cl.DecryptMessage(loadNetMsg.GetPayload().GetBody())
		if err != nil {
			t.Error(err)
			return
		}

		if pubKey.GetHasher().ToString() != cl.GetPubKey().GetHasher().ToString() {
			t.Error("load public key != init public key")
			return
		}

		pl := payload.LoadPayload64(decMsg)
		if pl.GetHead() != uint64(testutils.TcHead) {
			t.Error("load msg head != init head")
			return
		}

		if !bytes.Equal(pl.GetBody(), []byte(testutils.TcBody)) {
			t.Error("load msg body != init body")
			return
		}
	}
}

func newNetworkMessageWithData(cl client.IClient, networkKey, data string) (net_message.IMessage, error) {
	msg, err := cl.EncryptMessage(
		cl.GetPubKey(),
		payload.NewPayload64(uint64(testutils.TcHead), []byte(data)).ToBytes(),
	)
	if err != nil {
		return nil, err
	}
	netMsg := net_message.NewMessage(
		net_message.NewConstructSettings(&net_message.SConstructSettings{
			FSettings: net_message.NewSettings(&net_message.SSettings{
				FNetworkKey:   networkKey,
				FWorkSizeBits: testutils.TCWorkSize,
			}),
		}),
		payload.NewPayload32(0, msg),
	)
	return netMsg, nil
}

func newNetworkMessage(cl client.IClient, networkKey string) (net_message.IMessage, error) {
	return newNetworkMessageWithData(cl, networkKey, testutils.TcBody)
}
