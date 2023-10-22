package message

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/payload"
)

const (
	tcHead = 12345
	tcBody = "hello, world!"
)

func TestMessage(t *testing.T) {
	t.Parallel()

	pld := payload.NewPayload(tcHead, []byte(tcBody))
	msg := NewMessage(pld)

	if !bytes.Equal(msg.GetPayload().GetBody(), []byte(tcBody)) {
		t.Error("payload body not equal body in message")
		return
	}

	if !bytes.Equal(msg.GetHash(), getHash(pld.ToBytes())) {
		t.Error("payload hash not equal hash of message")
		return
	}

	if msg.GetPayload().GetHead() != tcHead {
		t.Error("payload head not equal head in message")
		return
	}

	msg1 := LoadMessage(msg.ToBytes())
	if !bytes.Equal(msg.GetPayload().ToBytes(), msg1.GetPayload().ToBytes()) {
		t.Error("load message not equal new message")
		return
	}

	if msg := LoadMessage([]byte{1}); msg != nil {
		t.Error("success load incorrect message")
		return
	}

	prng := random.NewStdPRNG()
	if msg := LoadMessage(prng.GetBytes(64)); msg != nil {
		t.Error("success load incorrect message")
		return
	}

	msgBytes := bytes.Join(
		[][]byte{
			{}, // pass payload
			getHash([]byte{}),
		},
		[]byte{},
	)
	if msg := LoadMessage(msgBytes); msg != nil {
		t.Error("success load incorrect payload")
		return
	}
}
