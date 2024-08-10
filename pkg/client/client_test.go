// nolint: goerr113
package client

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/payload/joiner"
	"github.com/number571/go-peer/pkg/utils"
	testutils "github.com/number571/go-peer/test/utils"
)

const (
	// tcMessageSize = (8 << 10)
	// tcKeySizeBits = 4096
	tcMessageSize = (2 << 10)
	tcKeySizeBits = 1024
)

var (
	// tgPrivKey = asymmetric.NewRSAPrivKey(tcKeySizeBits)
	tgPrivKey = asymmetric.LoadRSAPrivKey(testutils.Tc1PrivKey1024)

	tgMsgLimit = testNewClient().GetMessageLimit()
	tgMessages = []string{
		testutils.TcBody,
		"",
		"A",
		"AA",
		"AAA",
		"AAAA",
		"AAAAA",
		"AAAAAA",
		"AAAAAAA",
		"AAAAAAAA",
		"AAAAAAAAA",
		"AAAAAAAAAA",
		"AAAAAAAAAAA",
		"AAAAAAAAAAAA",
		"AAAAAAAAAAAAA",
		"AAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		random.NewCSPRNG().GetString(tgMsgLimit - encoding.CSizeUint64), // maximum size of message - payload64.head
	}
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SClientError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestClientPanicWithMessageSize(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()
	_ = NewClient(
		message.NewSettings(&message.SSettings{
			FMessageSizeBytes: 1024,
			FKeySizeBits:      512,
		}),
		tgPrivKey,
	)
}

func TestClientPanicWithKeySize(t *testing.T) {
	t.Parallel()

	testDiffKeySize(t)
	testLittleKeySize(t)
	testLittleMessageSize(t)
}

func testDiffKeySize(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()
	_ = NewClient(
		message.NewSettings(&message.SSettings{
			FMessageSizeBytes: tcMessageSize,
			FKeySizeBits:      4096,
		}),
		tgPrivKey,
	)
}

func testLittleKeySize(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()
	_ = NewClient(
		message.NewSettings(&message.SSettings{
			FMessageSizeBytes: tcMessageSize,
			FKeySizeBits:      128,
		}),
		asymmetric.NewRSAPrivKey(128),
	)
}

func testLittleMessageSize(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()
	_ = NewClient(
		message.NewSettings(&message.SSettings{
			FMessageSizeBytes: 128,
			FKeySizeBits:      tcKeySizeBits,
		}),
		tgPrivKey,
	)
}

func TestEncrypt(t *testing.T) {
	t.Parallel()

	client1 := testNewClient()
	client2 := testNewClient()

	pl := payload.NewPayload64(uint64(testutils.TcHead), []byte(testutils.TcBody))
	msg, err := client1.EncryptMessage(client2.GetPubKey(), pl.ToBytes())
	if err != nil {
		t.Error(err)
		return
	}

	// lMsg, _ := message.LoadMessage(client1.GetSettings(), msg)
	// os.WriteFile("test_binary.msg", lMsg.ToBytes(), 0644)
	// os.WriteFile("test_string.msg", []byte(lMsg.ToString()), 0644)

	_, decMsg, err := client2.DecryptMessage(msg)
	if err != nil {
		t.Error(err)
		return
	}

	decPl := payload.LoadPayload64(decMsg)
	if !bytes.Equal([]byte(testutils.TcBody), decPl.GetBody()) {
		t.Error("data not equal with decrypted data")
		return
	}
}

func TestDecrypt(t *testing.T) {
	t.Parallel()

	client1 := testNewClient()

	pl := payload.NewPayload64(uint64(testutils.TcHead), []byte(testutils.TcBody))
	msg1, err := client1.EncryptMessage(client1.GetPubKey(), pl.ToBytes())
	if err != nil {
		t.Error(err)
		return
	}

	if _, _, err := client1.DecryptMessage(msg1); err != nil {
		t.Error(err)
		return
	}

	msg, err := message.LoadMessage(client1.GetSettings(), msg1)
	if err != nil {
		t.Error(err)
		return
	}

	newEncd1 := make([]byte, len(msg.GetEncd()))
	copy(newEncd1, msg.GetEncd())
	newEncd1[0] ^= 1

	newMsg1 := message.NewMessage(msg.GetEnck(), newEncd1)
	if _, _, err := client1.DecryptMessage(newMsg1.ToBytes()); err == nil {
		t.Error("success decrypt invalid message")
		return
	}

	newEncd2 := make([]byte, len(msg.GetEncd()))
	copy(newEncd2, msg.GetEncd())
	newEncd2[symmetric.CAESBlockSize+8+1] ^= 1 // public key padding

	newMsg2 := message.NewMessage(msg.GetEnck(), newEncd2)
	if _, _, err := client1.DecryptMessage(newMsg2.ToBytes()); err == nil {
		t.Error("success decrypt invalid message (public key)")
		return
	}

	newEncd3 := make([]byte, len(msg.GetEncd()))
	copy(newEncd3, msg.GetEncd())
	newEncd3[symmetric.CAESBlockSize+196+1] ^= 1 // hash padding

	newMsg3 := message.NewMessage(msg.GetEnck(), newEncd3)
	if _, _, err := client1.DecryptMessage(newMsg3.ToBytes()); err == nil {
		t.Error("success decrypt invalid message (hash)")
		return
	}

	newEncd4 := make([]byte, len(msg.GetEncd()))
	copy(newEncd4, msg.GetEncd())
	newEncd4[symmetric.CAESBlockSize+236+1] ^= 1 // sign padding

	newMsg4 := message.NewMessage(msg.GetEnck(), newEncd4)
	if _, _, err := client1.DecryptMessage(newMsg4.ToBytes()); err == nil {
		t.Error("success decrypt invalid message (sign)")
		return
	}

	if _, _, err := client1.DecryptMessage(nil); err == nil {
		t.Error("success decrypt nil message")
		return
	}

	client1Ptr := client1.(*sClient)
	msg3, err := client1Ptr.tInvalidEncryptPayload1(client1.GetPubKey(), pl)
	if err != nil {
		t.Error(err)
		return
	}
	if _, _, err := client1.DecryptMessage(msg3.ToBytes()); err == nil {
		t.Error("success decrypt message with incorrect payload (1)")
		return
	}

	msg4, err := client1Ptr.tInvalidEncryptPayload2(client1.GetPubKey(), pl)
	if err != nil {
		t.Error(err)
		return
	}
	if _, _, err := client1.DecryptMessage(msg4.ToBytes()); err == nil {
		t.Error("success decrypt message with incorrect payload (2)")
		return
	}

	msg5, err := client1Ptr.tInvalidEncryptPayload3(client1.GetPubKey(), pl)
	if err != nil {
		t.Error(err)
		return
	}
	if _, _, err := client1.DecryptMessage(msg5.ToBytes()); err == nil {
		t.Error("success decrypt message with incorrect payload (3)")
		return
	}

	msg6, err := client1Ptr.tInvalidEncryptPayload4(client1.GetPubKey(), pl)
	if err != nil {
		t.Error(err)
		return
	}
	if _, _, err := client1.DecryptMessage(msg6.ToBytes()); err == nil {
		t.Error("success decrypt message with incorrect payload (4)")
		return
	}
}

func TestMessageSize(t *testing.T) {
	t.Parallel()

	client1 := testNewClient()

	for _, smsg := range tgMessages {
		pl := payload.NewPayload64(uint64(testutils.TcHead), []byte(smsg))
		msg, err := client1.EncryptMessage(client1.GetPubKey(), pl.ToBytes())
		if err != nil {
			t.Error(err)
			return
		}
		if uint64(len(msg)) != client1.GetSettings().GetMessageSizeBytes() {
			t.Error("got invalid message size bytes")
			return
		}
	}
}

func TestGetMessageLimit(t *testing.T) {
	t.Parallel()

	client1 := testNewClient()

	msg1 := random.NewCSPRNG().GetBytes(tgMsgLimit - encoding.CSizeUint64)
	pld1 := payload.NewPayload64(uint64(testutils.TcHead), msg1)
	if _, err := client1.EncryptMessage(client1.GetPubKey(), pld1.ToBytes()); err != nil {
		t.Error("message1 > message limit:", err)
		return
	}

	msg2 := random.NewCSPRNG().GetBytes(tgMsgLimit + 1)
	pld2 := payload.NewPayload64(uint64(testutils.TcHead), msg2)
	if _, err := client1.EncryptMessage(client1.GetPubKey(), pld2.ToBytes()); err == nil {
		t.Error("message2 > message limit but not alert:", err)
		return
	}
}

func testNewClient() IClient {
	return NewClient(
		message.NewSettings(&message.SSettings{
			FMessageSizeBytes: tcMessageSize,
			FKeySizeBits:      tcKeySizeBits,
		}),
		tgPrivKey,
	)
}

// INVALID_1

func (p *sClient) tInvalidEncryptPayload1(pRecv asymmetric.IPubKey, pPld payload.IPayload64) (message.IMessage, error) {
	var (
		msgLimitSize = p.GetMessageLimit()
		resultSize   = uint64(len(pPld.ToBytes()))
	)

	if resultSize > msgLimitSize {
		return nil, utils.MergeErrors(
			ErrLimitMessageSize,
			fmt.Errorf(
				"limit of message size without hex encoding = %d bytes < current payload size with additional padding = %d bytes",
				msgLimitSize,
				resultSize,
			),
		)
	}

	return p.tInvalidEncryptWithParams1(
		pRecv,
		pPld,
		msgLimitSize-resultSize,
	), nil
}

func (p *sClient) tInvalidEncryptWithParams1(pRecv asymmetric.IPubKey, pPld payload.IPayload64, pPadd uint64) message.IMessage {
	var (
		rand    = random.NewCSPRNG()
		salt    = rand.GetBytes(cSaltSize)
		session = rand.GetBytes(symmetric.CAESKeySize)
	)

	payloadBytes := pPld.ToBytes()
	doublePayload := payload.NewPayload64(
		uint64(len(payloadBytes))-1,
		bytes.Join(
			[][]byte{
				payloadBytes,
				rand.GetBytes(pPadd),
			},
			[]byte{},
		),
	)

	hash := hashing.NewHMACSHA256Hasher(salt, bytes.Join(
		[][]byte{
			p.GetPubKey().GetHasher().ToBytes(),
			pRecv.GetHasher().ToBytes(),
			doublePayload.ToBytes(),
		},
		[]byte{},
	)).ToBytes()

	encKey := pRecv.EncryptBytes(session)
	if encKey == nil {
		panic(ErrEncryptSymmetricKey)
	}

	cipher := symmetric.NewAESCipher(session)
	return message.NewMessage(
		encKey,
		cipher.EncryptBytes(joiner.NewBytesJoiner32([][]byte{
			p.GetPubKey().ToBytes(),
			salt,
			hash,
			p.fPrivKey.SignBytes(hash),
			doublePayload.ToBytes(),
		})),
	)
}

// INVALID_2

func (p *sClient) tInvalidEncryptPayload2(pRecv asymmetric.IPubKey, pPld payload.IPayload64) (message.IMessage, error) {
	var (
		msgLimitSize = p.GetMessageLimit()
		resultSize   = uint64(len(pPld.GetBody()))
	)

	if resultSize > msgLimitSize {
		return nil, ErrLimitMessageSize
	}

	return p.tInvalidEncryptWithParams2(
		pRecv,
		pPld,
		msgLimitSize-resultSize,
	)
}

func (p *sClient) tInvalidEncryptWithParams2(
	pRecv asymmetric.IPubKey,
	pPld payload.IPayload64,
	pPadd uint64,
) (message.IMessage, error) {
	var (
		rand    = random.NewCSPRNG()
		salt    = rand.GetBytes(cSaltSize)
		session = rand.GetBytes(symmetric.CAESKeySize)
	)

	data := joiner.NewBytesJoiner32([][]byte{
		bytes.Join(
			[][]byte{
				pPld.ToBytes(),
				rand.GetBytes(pPadd),
				{1, 2, 3, 4}, // uint32
			},
			[]byte{},
		),
	})

	hash := hashing.NewHMACSHA256Hasher(salt, bytes.Join(
		[][]byte{
			p.GetPubKey().GetHasher().ToBytes(),
			pRecv.GetHasher().ToBytes(),
			data,
		},
		[]byte{},
	)).ToBytes()

	encKey := pRecv.EncryptBytes(session)
	if encKey == nil {
		return nil, ErrEncryptSymmetricKey
	}

	cipher := symmetric.NewAESCipher(session)
	return message.NewMessage(
		encKey,
		cipher.EncryptBytes(joiner.NewBytesJoiner32([][]byte{
			p.GetPubKey().ToBytes(),
			salt,
			hash,
			p.fPrivKey.SignBytes(hash),
			data,
		})),
	), nil
}

// INVALID_3

func (p *sClient) tInvalidEncryptPayload3(pRecv asymmetric.IPubKey, pPld payload.IPayload64) (message.IMessage, error) {
	var (
		msgLimitSize = p.GetMessageLimit()
		resultSize   = uint64(len(pPld.GetBody()))
	)

	if resultSize > msgLimitSize {
		return nil, ErrLimitMessageSize
	}

	return p.tInvalidEncryptWithParams3(
		pRecv,
		pPld,
		msgLimitSize+encoding.CSizeUint64,
	)
}

func (p *sClient) tInvalidEncryptWithParams3(
	pRecv asymmetric.IPubKey,
	_ payload.IPayload64,
	pPadd uint64,
) (message.IMessage, error) {
	var (
		rand    = random.NewCSPRNG()
		salt    = rand.GetBytes(cSaltSize)
		session = rand.GetBytes(symmetric.CAESKeySize)
	)

	data := joiner.NewBytesJoiner32([][]byte{
		nil,
		rand.GetBytes(pPadd),
	})

	hash := hashing.NewHMACSHA256Hasher(salt, bytes.Join(
		[][]byte{
			p.GetPubKey().GetHasher().ToBytes(),
			pRecv.GetHasher().ToBytes(),
			data,
		},
		[]byte{},
	)).ToBytes()

	encKey := pRecv.EncryptBytes(session)
	if encKey == nil {
		return nil, ErrEncryptSymmetricKey
	}

	cipher := symmetric.NewAESCipher(session)
	return message.NewMessage(
		encKey,
		cipher.EncryptBytes(joiner.NewBytesJoiner32([][]byte{
			p.GetPubKey().ToBytes(),
			salt,
			hash,
			p.fPrivKey.SignBytes(hash),
			data,
		})),
	), nil
}

// INVALID_4

func (p *sClient) tInvalidEncryptPayload4(pRecv asymmetric.IPubKey, pPld payload.IPayload64) (message.IMessage, error) {
	var (
		msgLimitSize = p.GetMessageLimit()
		resultSize   = uint64(len(pPld.GetBody()))
	)

	if resultSize > msgLimitSize {
		return nil, ErrLimitMessageSize
	}

	return p.tInvalidEncryptWithParams4(
		pRecv,
		pPld,
		msgLimitSize-resultSize-encoding.CSizeUint64,
	)
}

func (p *sClient) tInvalidEncryptWithParams4(
	pRecv asymmetric.IPubKey,
	pPld payload.IPayload64,
	pPadd uint64,
) (message.IMessage, error) {
	var (
		rand    = random.NewCSPRNG()
		salt    = rand.GetBytes(cSaltSize)
		session = rand.GetBytes(symmetric.CAESKeySize)
	)

	data := joiner.NewBytesJoiner32([][]byte{
		pPld.ToBytes(),
		rand.GetBytes(pPadd),
	})

	hash := hashing.NewHMACSHA256Hasher(salt, bytes.Join(
		[][]byte{
			p.GetPubKey().GetHasher().ToBytes(),
			pRecv.GetHasher().ToBytes(),
			data,
		},
		[]byte{},
	)).ToBytes()

	hash[0] ^= 1

	encKey := pRecv.EncryptBytes(session)
	if encKey == nil {
		return nil, ErrEncryptSymmetricKey
	}

	cipher := symmetric.NewAESCipher(session)
	return message.NewMessage(
		encKey,
		cipher.EncryptBytes(joiner.NewBytesJoiner32([][]byte{
			p.GetPubKey().ToBytes(),
			salt,
			hash,
			p.fPrivKey.SignBytes(hash),
			data,
		})),
	), nil
}
