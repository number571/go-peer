package msgconv

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/client/message"
	testutils "github.com/number571/go-peer/test/_data"
)

func TestInvalidConvert(t *testing.T) {
	t.Parallel()

	if res := FromBytesToString([]byte{123}); res != "" {
		t.Error("success convert invalid bytes to string")
		return
	}
	if res := FromStringToBytes("123"); res != nil {
		t.Error("success convert invalid string to bytes (split)")
		return
	}
	if res := FromStringToBytes("123" + message.CSeparator + "!@#"); res != nil {
		t.Error("success convert invalid string to bytes (hex decode)")
		return
	}
}

func TestConvert(t *testing.T) {
	t.Parallel()

	params := message.NewSettings(&message.SSettings{
		FMessageSizeBytes: (2 << 10),
		FWorkSizeBits:     testutils.TCWorkSize,
	})

	msg1 := message.LoadMessage(params, FromBytesToString(testutils.TCBinaryMessage))
	if msg1 == nil {
		t.Error("fromBytesToString result is invalid")
		return
	}
	if !bytes.Equal(msg1.ToBytes(), testutils.TCBinaryMessage) {
		t.Error("msg1 bytes not equal with original")
		return
	}

	msg2 := message.LoadMessage(params, FromStringToBytes(testutils.TCStringMessage))
	if msg2 == nil {
		t.Error("fromStringToBytes result is invalid")
		return
	}
	if msg2.ToString() != testutils.TCStringMessage {
		t.Error("msg2 string not equal with original")
		return
	}
}
