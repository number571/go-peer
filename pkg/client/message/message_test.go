package message

import (
	"bytes"
	"os"
	"testing"
)

func TestMessage(t *testing.T) {
	msgBytes, err := os.ReadFile("message.json")
	if err != nil {
		t.Error(err)
		return
	}

	msg := LoadMessage(msgBytes, 100<<10, 10)
	if msg == nil {
		t.Error("failed load message")
		return
	}

	if !bytes.Equal(msg.Bytes(), msg.Bytes()) {
		t.Error("load message not equal new message")
		return
	}
}
