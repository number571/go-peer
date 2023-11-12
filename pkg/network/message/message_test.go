package message

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload"
	testutils "github.com/number571/go-peer/test/_data"
)

const (
	tcHead       = 12345
	tcBody       = "hello, world!"
	tcNetworkKey = "network_key"
	tcProof      = 1096
)

func TestMessage(t *testing.T) {
	t.Parallel()

	pld := payload.NewPayload(tcHead, []byte(tcBody))
	sett := NewSettings(&SSettings{
		FWorkSizeBits: testutils.TCWorkSize,
		FNetworkKey:   tcNetworkKey,
	})
	msg := NewMessage(sett, pld)

	if !bytes.Equal(msg.GetPayload().GetBody(), []byte(tcBody)) {
		t.Error("payload body not equal body in message")
		return
	}

	if !bytes.Equal(msg.GetHash(), getHash(tcNetworkKey, pld.ToBytes())) {
		t.Error("payload hash not equal hash of message")
		return
	}

	if msg.GetPayload().GetHead() != tcHead {
		t.Error("payload head not equal head in message")
		return
	}

	if msg.GetProof() != tcProof {
		t.Error("got invalid proof")
		return
	}

	msg1 := LoadMessage(sett, msg.ToBytes())
	if !bytes.Equal(msg.GetPayload().ToBytes(), msg1.GetPayload().ToBytes()) {
		t.Error("load message not equal new message")
		return
	}

	if msg := LoadMessage(sett, []byte{1}); msg != nil {
		t.Error("success load incorrect message")
		return
	}

	randBytes := random.NewStdPRNG().GetBytes(encoding.CSizeUint64 + hashing.CSHA256Size)
	if msg := LoadMessage(sett, randBytes); msg != nil {
		t.Error("success load incorrect message")
		return
	}

	prng := random.NewStdPRNG()
	if msg := LoadMessage(sett, prng.GetBytes(64)); msg != nil {
		t.Error("success load incorrect message")
		return
	}

	msgBytes := bytes.Join(
		[][]byte{
			{}, // pass payload
			getHash("_", []byte{}),
		},
		[]byte{},
	)
	if msg := LoadMessage(sett, msgBytes); msg != nil {
		t.Error("success load incorrect payload")
		return
	}
}
