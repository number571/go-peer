package macro

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/payload/joiner"

	_ "embed"
)

var (
	//go:embed test_binary.msg
	tgBinaryMessage []byte

	//go:embed test_string.msg
	tgStringMessage string
)

func TestPanicnewMessage(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()

	_ = newMessage([]byte{}, []byte{})
}

func TestInvalidMessage(t *testing.T) {
	t.Parallel()

	msgSize := uint64(2 << 10)

	if _, err := loadMessage(msgSize, struct{}{}); err == nil {
		t.Error("success load message with unknown type")
		return
	}

	if _, err := loadMessage(msgSize, []byte{123}); err == nil {
		t.Error("success load invalid message")
		return
	}

	msgBytes := joiner.NewBytesJoiner32([][]byte{[]byte("aaa"), []byte("bbb")})
	if _, err := loadMessage(msgSize, msgBytes); err == nil {
		t.Error("success load invalid message")
		return
	}

	if _, err := loadMessage(1, msgBytes); err == nil {
		t.Error("success load message with keysize > msgsize")
		return
	}
}

func TestMessage(t *testing.T) {
	t.Parallel()

	msgSize := uint64(8 << 10)

	msg1, err := loadMessage(msgSize, tgBinaryMessage)
	if err != nil {
		t.Error(err)
		return
	}
	testMessage(t, msgSize, msg1)

	msg2, err := loadMessage(msgSize, tgStringMessage)
	if err != nil {
		t.Error(err)
		return
	}
	testMessage(t, msgSize, msg2)
}

func testMessage(t *testing.T, msgSize uint64, msg IMessage) {
	if !bytes.Equal(msg.ToBytes(), tgBinaryMessage) {
		t.Error("invalid convert to bytes")
		return
	}

	if msg.ToString() != tgStringMessage {
		t.Error("invalid convert to string")
		return
	}

	msgBytes := bytes.Join([][]byte{msg.GetEnck(), msg.GetEncd()}, []byte{})
	if _, err := loadMessage(msgSize, msgBytes); err != nil {
		t.Error("new message is invalid")
		return
	}
}
