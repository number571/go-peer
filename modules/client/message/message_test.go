package message

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/modules/payload"
)

const (
	tcHead = 12345
	tcBody = "hello, world!"
)

func TestMessage(t *testing.T) {
	msg := testNewMessage()
	msg1 := LoadMessage(msg.Bytes())

	if !bytes.Equal(msg.Bytes(), msg1.Bytes()) {
		t.Error("load message not equal new message")
		return
	}
}

func testNewMessage() IMessage {
	return &SMessage{
		FHead: SHeadMessage{
			FSession: "session-key",
		},
		FBody: SBodyMessage{
			FPayload: string(payload.NewPayload(tcHead, []byte(tcBody)).Bytes()),
		},
	}
}
