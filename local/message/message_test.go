package message

import (
	"bytes"
	"testing"
)

const (
	tcTitle = "title"
	tcData  = "data"
)

func TestMessage(t *testing.T) {
	msg := NewMessage([]byte(tcTitle), []byte(tcData))

	packBytes := msg.ToPackage().Bytes()
	decMsg := LoadPackage(packBytes).ToMessage()
	if decMsg == nil {
		t.Error("decode message is nil")
		return
	}

	if !bytes.Equal(msg.Body().Data(), decMsg.Body().Data()) {
		t.Error("data not equal with decoded version of message")
		return
	}
}
