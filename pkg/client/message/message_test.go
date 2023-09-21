package message

import (
	"bytes"
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/encoding"
	testutils "github.com/number571/go-peer/test/_data"
)

const (
	tcSalt    = "2823c693930a2d230415368221c36f5c1c77accd82c949e5de13bbe8ef00275b72d182d79759b1aef582649c99afd409"
	tcSession = "21273151593d606632e0015107c1f204af68dd1a4977a6e7717260936dc4043e1f57410914f0dcffe372d4fc948ce0b12e9449e9e35556416d6e0f4f1d6c08807b0e98588aeea85c3f4a62512d82d159734d80ea16af0380a0d77f57c5444eaaf4cd9eed67f13fbc4699c608736286cba3e962fdc1642beac195ec59dd2d926a"
	tcSender  = "bc2ebe882b76fd7ab77ec19416328902f685dfedec57285d60ac154d7fa67a9d735bd3f22ba4f7394fead3c35bbcc30b8f67dccd15c8a663f4014785124eaf098ed83027870949ed4bc0a0516a9ac901f185b71d7765e14b00edf2743cd02b266695977c8730f891d76c9a69795806edb61426b8bc2f3d32de0a074995844eed93a5648532ff0ba8ee5b1ebd6bd3b33620b7d938852776b8b99276fa"

	tcSign  = "fb3ce5d111c74a0b8638f24f8ff200f64ca0e88cda1fd483783930b08e465fa9fc9565a0a3afbdfdf3f463bc77e526f2c41c6ddd2dae5d6f90e741442e2939731cbdad4071c29eff83dff932589b2cbfd8fa8a5fac19de4c40c3adde4cde1235c0bbf053b0e04e826993f8060a50c671c6bf56ce24fe4e921b60f6ca2239932ebd1b8c8556d5a2ac13e5ef1d8ea9cca8"
	tcHash  = "e1cc8da32433c2d7f99b38042d7b73db291bd803c55f3c83745ae3ebae6baffe"
	tcProof = 4512
)

func TestConvert(t *testing.T) {
	params := NewSettings(&SSettings{
		FMessageSizeBytes: (2 << 10),
		FWorkSizeBits:     testutils.TCWorkSize,
	})

	msgBytes, err := os.ReadFile("test_binary.msg")
	if err != nil {
		t.Error(err)
		return
	}
	msg1 := LoadMessage(params, FromBytesToString(msgBytes))
	if msg1 == nil {
		t.Error("fromBytesToString result is invalid")
		return
	}
	if !bytes.Equal(msg1.ToBytes(), msgBytes) {
		t.Error("msg1 bytes not equal with original")
		return
	}

	msgStrBytes, err := os.ReadFile("test_string.msg")
	if err != nil {
		t.Error(err)
		return
	}
	msg2 := LoadMessage(params, FromStringToBytes(string(msgStrBytes)))
	if msg2 == nil {
		t.Error("fromStringToBytes result is invalid")
		return
	}
	if msg2.ToString() != string(msgStrBytes) {
		t.Error("msg2 string not equal with original")
		return
	}
}

func TestMessage(t *testing.T) {
	params := NewSettings(&SSettings{
		FMessageSizeBytes: (2 << 10),
		FWorkSizeBits:     testutils.TCWorkSize,
	})

	msgBytes, err := os.ReadFile("test_binary.msg")
	if err != nil {
		t.Error(err)
		return
	}
	msg1 := LoadMessage(params, msgBytes)
	if msg1 == nil {
		t.Error("failed load message")
		return
	}
	testMessage(t, msg1)

	msgStrBytes, err := os.ReadFile("test_string.msg")
	if err != nil {
		t.Error(err)
		return
	}
	msg2 := LoadMessage(params, string(msgStrBytes))
	if msg2 == nil {
		t.Error("failed load message")
		return
	}
	testMessage(t, msg2)
}

func testMessage(t *testing.T, msg IMessage) {
	if !bytes.Equal(msg.GetHead().GetSalt(), encoding.HexDecode(tcSalt)) {
		t.Error("incorrect salt value")
		return
	}

	if !bytes.Equal(msg.GetHead().GetSession(), encoding.HexDecode(tcSession)) {
		t.Error("incorrect session value")
		return
	}

	if !bytes.Equal(msg.GetHead().GetSender(), encoding.HexDecode(tcSender)) {
		t.Error("incorrect sender value")
		return
	}

	if !bytes.Equal(msg.GetBody().GetSign(), encoding.HexDecode(tcSign)) {
		t.Error("incorrect sign value")
		return
	}

	if !bytes.Equal(msg.GetBody().GetHash(), encoding.HexDecode(tcHash)) {
		t.Error("incorrect hash value")
		return
	}

	if msg.GetBody().GetProof() != tcProof {
		t.Error("incorrect proof value")
		return
	}

	// pass payload -> message is encrypted
}
