package message

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload"
	testutils "github.com/number571/go-peer/test/utils"
)

var (
	_ payload.IPayload = &sInvalidPayload{}
)

const (
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

func TestMessage(t *testing.T) {
	t.Parallel()

	pld := payload.NewPayload(tcHead, []byte(tcBody))
	sett := NewSettings(&SSettings{
		FWorkSizeBits: testutils.TCWorkSize,
		FNetworkKey:   tcNetworkKey,
	})
	msg := NewMessage(sett, pld, 1, 0)

	if !bytes.Equal(msg.GetPayload().GetBody(), []byte(tcBody)) {
		t.Error("payload body not equal body in message")
		return
	}

	// if !bytes.Equal(msg.GetHash(), getAuthHash(tcNetworkKey, msg.GetSalt()[0], pld.ToBytes())) {
	// 	t.Error("payload hash not equal hash of message")
	// 	return
	// }

	if msg.GetPayload().GetHead() != tcHead {
		t.Error("payload head not equal head in message")
		return
	}

	for i := 0; i < 10; i++ {
		msgN := NewMessage(sett, pld, 1, 0)
		if msgN.GetProof() == 0 {
			continue
		}
		msgL, err := LoadMessage(sett, msgN.ToBytes())
		if err != nil {
			t.Error(err)
			return
		}
		if msgN.GetProof() != msgL.GetProof() {
			t.Error("got invalid proof")
			return
		}
		break
	}

	msg1, err := LoadMessage(sett, msg.ToBytes())
	if err != nil {
		t.Error(err)
		return
	}
	if !bytes.Equal(msg.GetPayload().ToBytes(), msg1.GetPayload().ToBytes()) {
		t.Error("load message not equal new message")
		return
	}

	msg2, err := LoadMessage(sett, msg.ToString())
	if err != nil {
		t.Error(err)
		return
	}
	if !bytes.Equal(msg.GetPayload().ToBytes(), msg2.GetPayload().ToBytes()) {
		t.Error("load message not equal new message")
		return
	}

	// msg3 := NewMessage(sett, pld, 1).(*sMessage)
	// msg3.fHash = random.NewStdPRNG().GetBytes(hashing.CSHA256Size)
	// msg3.fProof = puzzle.NewPoWPuzzle(testutils.TCWorkSize).ProofBytes(msg3.fHash, 1)
	// if _, err := LoadMessage(sett, msg3.ToBytes()); err == nil {
	// 	t.Error("success load with invalid hash")
	// 	return
	// }

	msg4 := NewMessage(sett, &sInvalidPayload{}, 1, 0)
	if _, err := LoadMessage(sett, msg4.ToBytes()); err == nil {
		t.Error("success load with invalid payload")
		return
	}

	if _, err := LoadMessage(sett, struct{}{}); err == nil {
		t.Error("success load with unknown type of message")
		return
	}

	if _, err := LoadMessage(sett, []byte{1}); err == nil {
		t.Error("success load incorrect message")
		return
	}

	randBytes := random.NewStdPRNG().GetBytes(encoding.CSizeUint64 + hashing.CSHA256Size)
	if _, err := LoadMessage(sett, randBytes); err == nil {
		t.Error("success load incorrect message")
		return
	}

	prng := random.NewStdPRNG()
	if _, err := LoadMessage(sett, prng.GetBytes(64)); err == nil {
		t.Error("success load incorrect message")
		return
	}

	msgBytes := bytes.Join(
		[][]byte{
			{}, // pass payload
			getAuthHash("_", []byte{}, []byte{}),
		},
		[]byte{},
	)
	if _, err := LoadMessage(sett, msgBytes); err == nil {
		t.Error("success load incorrect payload")
		return
	}
}
