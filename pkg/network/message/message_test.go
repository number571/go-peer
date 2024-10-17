package message

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/keybuilder"
	"github.com/number571/go-peer/pkg/crypto/puzzle"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/payload/joiner"
	testutils "github.com/number571/go-peer/test/utils"
)

var (
	_ payload.IPayload64 = &sInvalidPayload{}
)

const (
	tcLimitVoid  = 128
	tcHead       = 12345
	tcBody       = "hello, world!"
	tcNetworkKey = "network_key_1"
)

type sInvalidPayload struct{}

func (p *sInvalidPayload) GetHead() uint64 {
	return 1
}

func (p *sInvalidPayload) GetBody() []byte {
	return []byte{}
}

func (p *sInvalidPayload) ToBytes() []byte {
	return []byte{123}
}

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SMessageError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestSettings(t *testing.T) {
	t.Parallel()

	for i := 0; i < 1; i++ {
		testSettings(t, i)
	}
}

func testSettings(t *testing.T, n int) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()
	switch n { // nolint: gocritic
	case 0:
		_ = NewConstructSettings(&SConstructSettings{})
	}
}

func TestMessage(t *testing.T) {
	t.Parallel()

	pld := payload.NewPayload32(tcHead, []byte(tcBody))
	sett := NewConstructSettings(&SConstructSettings{
		FSettings: NewSettings(&SSettings{
			FWorkSizeBits: testutils.TCWorkSize,
			FNetworkKey:   tcNetworkKey,
		}),
		FRandMessageSizeBytes: tcLimitVoid,
	})

	msgTmp := NewMessage(sett, pld)
	if !bytes.Equal(msgTmp.GetPayload().GetBody(), []byte(tcBody)) {
		t.Error("payload body not equal body in message")
		return
	}

	msg, err := LoadMessage(sett.GetSettings(), msgTmp.ToBytes())
	if err != nil {
		t.Error(err)
		return
	}

	if msgTmp.ToString() != msg.ToString() {
		t.Error("msgTmp != msg")
		return
	}

	voidBytes := msg.GetRand()
	if len(voidBytes) > tcLimitVoid {
		t.Error("got length void bytes > limit")
		return
	}

	if (len(msg.ToBytes()) - len(voidBytes)) != len(msg.GetPayload().GetBody())+CMessageHeadSize {
		t.Error("invalid message size in bytes with void")
		return
	}

	payloadSize := encoding.Uint32ToBytes(uint32(len(pld.ToBytes())))
	voidSize := encoding.Uint32ToBytes(uint32(len(voidBytes)))
	payloadRandBytes := bytes.Join(
		[][]byte{
			payloadSize[:],
			pld.ToBytes(),
			voidSize[:],
			voidBytes,
		},
		[]byte{},
	)

	keyBuilder := keybuilder.NewKeyBuilder(0, []byte{}) // the network_key must have good entropy
	key := keyBuilder.Build(tcNetworkKey, symmetric.CCipherKeySize)
	newHash := hashing.NewHMACHasher(key, payloadRandBytes).ToBytes()
	if !bytes.Equal(msg.GetHash(), newHash) {
		t.Error("payload hash not equal hash of message")
		return
	}

	if msg.GetPayload().GetHead() != tcHead {
		t.Error("payload head not equal head in message")
		return
	}

	newSett := NewConstructSettings(&SConstructSettings{
		FSettings: NewSettings(&SSettings{
			FWorkSizeBits: testutils.TCWorkSize,
			FNetworkKey:   tcNetworkKey,
		}),
	})

	for i := 0; i < 10; i++ {
		msgN := NewMessage(newSett, pld)
		if msgN.GetProof() == 0 {
			continue
		}
		msgL, err := LoadMessage(newSett.GetSettings(), msgN.ToBytes())
		if err != nil {
			t.Error(err)
			return
		}
		if msgN.GetProof() != msgL.GetProof() {
			t.Error("got invalid proof")
			return
		}
		if len(msgN.ToBytes()) != len(msgL.ToBytes()) {
			t.Error("new msg size != load msg size")
			return
		}
		if len(msgN.ToBytes()) != CMessageHeadSize+len(pld.GetBody()) {
			t.Error("msg size != head size + payload body")
			return
		}
		break
	}

	msg1, err := LoadMessage(sett.GetSettings(), msg.ToBytes())
	if err != nil {
		t.Error(err)
		return
	}
	if !bytes.Equal(msg.GetPayload().ToBytes(), msg1.GetPayload().ToBytes()) {
		t.Error("load message not equal new message")
		return
	}

	msg2, err := LoadMessage(sett.GetSettings(), msg.ToString())
	if err != nil {
		t.Error(err)
		return
	}
	if !bytes.Equal(msg.GetPayload().ToBytes(), msg2.GetPayload().ToBytes()) {
		t.Error("load message not equal new message")
		return
	}

	msg3 := NewMessage(sett, pld).(*sMessage)
	msg3.fEncd[0] ^= 1
	if _, err := LoadMessage(sett.GetSettings(), msg3.ToBytes()); err == nil {
		t.Error("success load with invalid encd")
		return
	}

	if _, err := LoadMessage(sett.GetSettings(), struct{}{}); err == nil {
		t.Error("success load with unknown type of message")
		return
	}

	if _, err := LoadMessage(sett.GetSettings(), []byte{1}); err == nil {
		t.Error("success load incorrect message")
		return
	}

	if _, err := LoadMessage(sett.GetSettings(), []byte{1}); err == nil {
		t.Error("success load incorrect message")
		return
	}

	randBytes := random.NewRandom().GetBytes(encoding.CSizeUint64 + hashing.CHasherSize)
	if _, err := LoadMessage(sett.GetSettings(), randBytes); err == nil {
		t.Error("success load incorrect message")
		return
	}

	prng := random.NewRandom()
	if _, err := LoadMessage(sett.GetSettings(), prng.GetBytes(64)); err == nil {
		t.Error("success load incorrect message")
		return
	}

	msgBytes := bytes.Join(
		[][]byte{
			{}, // pass payload
			hashing.NewHasher([]byte{}).ToBytes(),
		},
		[]byte{},
	)
	if _, err := LoadMessage(sett.GetSettings(), msgBytes); err == nil {
		t.Error("success load incorrect payload")
		return
	}

	if _, err := LoadMessage(sett.GetSettings(), tNewInvalidMessage1(sett, pld).ToBytes()); err == nil {
		t.Error("success load invalid message 1")
		return
	}

	if _, err := LoadMessage(sett.GetSettings(), tNewInvalidMessage2(sett, pld).ToBytes()); err == nil {
		t.Error("success load invalid message 2")
		return
	}

	if _, err := LoadMessage(sett.GetSettings(), tNewInvalidMessage3(sett, pld).ToBytes()); err == nil {
		t.Error("success load invalid message 3")
		return
	}
}

func tNewInvalidMessage1(pSett IConstructSettings, pPld payload.IPayload32) IMessage {
	bytesJoiner := joiner.NewBytesJoiner32([][]byte{pPld.ToBytes()})
	sett := pSett.GetSettings()

	keyBuilder := keybuilder.NewKeyBuilder(0, []byte{}) // the network_key must have good entropy
	key := keyBuilder.Build(tcNetworkKey, symmetric.CCipherKeySize)
	hash := hashing.NewHMACHasher(key, bytesJoiner).ToBytes()

	proof := puzzle.NewPoWPuzzle(sett.GetWorkSizeBits()).ProofBytes(hash, pSett.GetParallel())
	proofBytes := encoding.Uint64ToBytes(proof)

	cipher := symmetric.NewCipher(key)
	return &sMessage{
		fEncd: cipher.EncryptBytes(bytes.Join(
			[][]byte{
				proofBytes[:],
				hash,
				bytesJoiner,
			},
			[]byte{},
		)),
		fHash:    hash,
		fProof:   proof,
		fPayload: pPld,
	}
}

func tNewInvalidMessage2(pSett IConstructSettings, pPld payload.IPayload32) IMessage {
	prng := random.NewRandom()
	sett := pSett.GetSettings()

	voidBytes := prng.GetBytes(prng.GetUint64() % (pSett.GetRandMessageSizeBytes() + 1))
	bytesJoiner := joiner.NewBytesJoiner32([][]byte{pPld.ToBytes(), voidBytes})

	keyBuilder := keybuilder.NewKeyBuilder(0, []byte{}) // the network_key must have good entropy
	key := keyBuilder.Build(tcNetworkKey, symmetric.CCipherKeySize)
	hash := hashing.NewHMACHasher(key, bytesJoiner).ToBytes()

	hash[0] ^= 1

	proof := puzzle.NewPoWPuzzle(sett.GetWorkSizeBits()).ProofBytes(hash, pSett.GetParallel())
	proofBytes := encoding.Uint64ToBytes(proof)

	cipher := symmetric.NewCipher(key)
	return &sMessage{
		fEncd: cipher.EncryptBytes(bytes.Join(
			[][]byte{
				proofBytes[:],
				hash,
				bytesJoiner,
			},
			[]byte{},
		)),
		fHash:    hash,
		fRand:    voidBytes,
		fProof:   proof,
		fPayload: pPld,
	}
}

func tNewInvalidMessage3(pSett IConstructSettings, pPld payload.IPayload32) IMessage {
	prng := random.NewRandom()
	sett := pSett.GetSettings()

	voidBytes := prng.GetBytes(prng.GetUint64() % (pSett.GetRandMessageSizeBytes() + 1))
	bytesJoiner := joiner.NewBytesJoiner32([][]byte{nil, voidBytes})

	keyBuilder := keybuilder.NewKeyBuilder(0, []byte{}) // the network_key must have good entropy
	key := keyBuilder.Build(tcNetworkKey, symmetric.CCipherKeySize)
	hash := hashing.NewHMACHasher(key, bytesJoiner).ToBytes()

	proof := puzzle.NewPoWPuzzle(sett.GetWorkSizeBits()).ProofBytes(hash, pSett.GetParallel())
	proofBytes := encoding.Uint64ToBytes(proof)

	cipher := symmetric.NewCipher(key)
	return &sMessage{
		fEncd: cipher.EncryptBytes(bytes.Join(
			[][]byte{
				proofBytes[:],
				hash,
				bytesJoiner,
			},
			[]byte{},
		)),
		fHash:    hash,
		fRand:    voidBytes,
		fProof:   proof,
		fPayload: pPld,
	}
}
