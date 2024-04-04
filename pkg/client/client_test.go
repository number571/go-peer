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
	"github.com/number571/go-peer/pkg/utils"
	testutils "github.com/number571/go-peer/test/utils"
)

const (
	tcPubKey1023Bit = "PubKey{30818802818067D7C7F48CCCF318A504721D8521ED04DFD24DE947D50476212E279DADA7F627102D6140B922BC30E778BE7FD120E73D02C407E4D053D55C8F3CFACED11CCB8DC96539F51E34B32544EAB8812129BE7AD107E9A352E70F912962D6D0BABAB629F20332A3FBA66BC28D2944F44859959788428300F58DA0873A6796073B17CEEB0203010001}"
	tcMessageSize   = (2 << 10)
	tcKeySizeBits   = 1024
)

var (
	tgMsgLimit = testNewClient().GetMessageLimit()
	tgPrivKey  = asymmetric.LoadRSAPrivKey(testutils.Tc1PrivKey1024)
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
		random.NewStdPRNG().GetString(tgMsgLimit), // maximum size of message
	}
)

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
			FKeySizeBits:      tcKeySizeBits,
		}),
		tgPrivKey,
	)
}

func TestClientPanicWithKeySize(t *testing.T) {
	t.Parallel()

	testDiffKeySize(t)
	testLittleKeySize(t)
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

func TestEncrypt(t *testing.T) {
	t.Parallel()

	client1 := testNewClient()
	client2 := testNewClient()

	_ = client1.GetSettings()
	_ = client1.GetPrivKey()

	pl := payload.NewPayload(uint64(testutils.TcHead), []byte(testutils.TcBody))
	msg, err := client1.EncryptPayload(client2.GetPubKey(), pl)
	if err != nil {
		t.Error(err)
		return
	}

	// os.WriteFile("test_binary.msg", msg.ToBytes(), 0644)
	// os.WriteFile("test_string.msg", []byte(msg.ToString()), 0644)

	_, decPl, err := client2.DecryptMessage(msg)
	if err != nil {
		t.Error(err)
		return
	}

	if !bytes.Equal([]byte(testutils.TcBody), decPl.GetBody()) {
		t.Error("data not equal with decrypted data")
		return
	}
}

func TestDecrypt(t *testing.T) {
	t.Parallel()

	client1 := testNewClient()

	pl := payload.NewPayload(uint64(testutils.TcHead), []byte(testutils.TcBody))
	msg, err := client1.EncryptPayload(client1.GetPubKey(), pl)
	if err != nil {
		t.Error(err)
		return
	}

	if _, _, err := client1.DecryptMessage(nil); err == nil {
		t.Error("success decrypt nil message")
		return
	}

	sMsg := msg.(*message.SMessage)

	sMsg1 := *sMsg
	sMsg1.FHash = "0000000000000000" + "0000000000000000000000000000000000000000000000000000000000000000"
	if _, _, err := client1.DecryptMessage(&sMsg1); err == nil {
		t.Error("success decrypt message with incorrect hash")
		return
	}

	sMsg3 := *sMsg
	sMsg3.FEncKey = "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
	if _, _, err := client1.DecryptMessage(&sMsg3); err == nil {
		t.Error("success decrypt message with incorrect session key")
		return
	}

	sMsg4 := *sMsg
	sMsg4.FPubKey = "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
	if _, _, err := client1.DecryptMessage(&sMsg4); err == nil {
		t.Error("success decrypt message with incorrect sender key (incorrect)")
		return
	}

	somePubKey := asymmetric.LoadRSAPubKey(tcPubKey1023Bit)
	msg2, err := client1.EncryptPayload(somePubKey, pl)
	if err != nil {
		t.Error(err)
		return
	}
	if _, _, err := client1.DecryptMessage(msg2); err == nil {
		t.Error("success decrypt message with incorrect sender key (size)")
		return
	}

	sMsg6 := *sMsg
	sMsg6.FPayload = []byte("111")
	if _, _, err := client1.DecryptMessage(&sMsg6); err == nil {
		t.Error("success decrypt message with incorrect payload (iv block)")
		return
	}

	sMsg7 := *sMsg
	sMsg7.FPayload = random.NewStdPRNG().GetBytes(948)
	if _, _, err := client1.DecryptMessage(&sMsg7); err == nil {
		t.Error("success decrypt message with incorrect payload (payload is nil)")
		return
	}

	sMsg8 := *sMsg
	sMsg8.FSalt = "00000000000000000000000000000000" + "0000000000000000000000000000000000000000000000000000000000000000"
	if _, _, err := client1.DecryptMessage(&sMsg8); err == nil {
		t.Error("success decrypt message with incorrect salt")
		return
	}

	sMsg9 := *sMsg
	sMsg9.FSign = "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
	if _, _, err := client1.DecryptMessage(&sMsg9); err == nil {
		t.Error("success decrypt message with incorrect sign")
		return
	}

	client1Ptr := client1.(*sClient)
	msg3, err := client1Ptr.tInvalidEncryptPayload(client1.GetPubKey(), pl)
	if err != nil {
		t.Error(err)
		return
	}
	if _, _, err := client1.DecryptMessage(msg3); err == nil {
		t.Error("success decrypt message with incorrect payload (double)")
		return
	}
}

func TestMessageSize(t *testing.T) {
	t.Parallel()

	client1 := testNewClient()
	sizes := make([]int, 0, len(tgMessages))

	for _, smsg := range tgMessages {
		pl := payload.NewPayload(uint64(testutils.TcHead), []byte(smsg))
		msg, err := client1.EncryptPayload(client1.GetPubKey(), pl)
		if err != nil {
			t.Error(err)
			return
		}
		sizes = append(sizes, len(msg.ToBytes()))
	}

	for i := 0; i < len(sizes)-1; i++ {
		if sizes[i] != sizes[i+1] {
			t.Errorf(
				"len bytes of different messages = id(%d, %d) not equals = size(%d, %d)",
				i, i+1,
				sizes[i], sizes[i+1],
			)
			return
		}
	}
}

func TestGetMessageLimit(t *testing.T) {
	t.Parallel()

	client1 := testNewClient()

	msg1 := random.NewStdPRNG().GetBytes(tgMsgLimit)
	pld1 := payload.NewPayload(uint64(testutils.TcHead), []byte(msg1))
	if _, err := client1.EncryptPayload(client1.GetPubKey(), pld1); err != nil {
		t.Error("message1 > message limit:", err)
		return
	}

	msg2 := random.NewStdPRNG().GetBytes(tgMsgLimit + 1)
	pld2 := payload.NewPayload(uint64(testutils.TcHead), []byte(msg2))
	if _, err := client1.EncryptPayload(client1.GetPubKey(), pld2); err == nil {
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

func (p *sClient) tInvalidEncryptPayload(pRecv asymmetric.IPubKey, pPld payload.IPayload) (message.IMessage, error) {
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

	return p.tInvalidEncryptWithParams(
		pRecv,
		pPld,
		msgLimitSize-resultSize,
	), nil
}

func (p *sClient) tInvalidEncryptWithParams(pRecv asymmetric.IPubKey, pPld payload.IPayload, pPadd uint64) message.IMessage {
	var (
		rand    = random.NewStdPRNG()
		salt    = rand.GetBytes(symmetric.CAESKeySize)
		session = rand.GetBytes(symmetric.CAESKeySize)
	)

	payloadBytes := pPld.ToBytes()
	doublePayload := payload.NewPayload(
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

	cipher := symmetric.NewAESCipher(session)
	return &message.SMessage{
		FPubKey:  encoding.HexEncode(cipher.EncryptBytes(p.GetPubKey().ToBytes())),
		FEncKey:  encoding.HexEncode(pRecv.EncryptBytes(session)),
		FSalt:    encoding.HexEncode(cipher.EncryptBytes(salt)),
		FHash:    encoding.HexEncode(cipher.EncryptBytes(hash)),
		FSign:    encoding.HexEncode(cipher.EncryptBytes(p.fPrivKey.SignBytes(hash))),
		FPayload: cipher.EncryptBytes(doublePayload.ToBytes()), // JSON field to raw Body (no need HEX encode)
	}
}
