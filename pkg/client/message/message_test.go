package message

import (
	"bytes"
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/encoding"
	testutils "github.com/number571/go-peer/test/_data"
)

const (
	tcSalt    = "f82291dde6c82aefad9c56462a5537a488be28a9decc1b5002c71e585a678283bad09298134d83bed2ad9fdf0a84a549"
	tcSession = "0eaf199a0458c25ba7433b06a1ca58cb9938e3a29208aff0a1aa7e7b52beced29cce2e309e89d09321bdee9feea9b084bf9aab4e361a6d15f2d502e451271ea8e2e87c6da329c9019e42356fb0e50afb47aacb15e0586c3d1a678d579f82e636ecb6cb3e8b8bcc082621e10952a5f4619d539bbcb9639a5ad546549490a0f9de"
	tcSender  = "ced82379a18f3fe84d839f8b4d8265c671a861b5140dc9021398c621d95c60365f8e7b19cdc47af1c0a92aabfc80c98731af5de901dbedb5d6d52b55c8bb3a7dab411a46891fb07f4308374139d55da01a952d637f00d8e3fbb457e504d6724e8b0b33563b0f26998434907eb6818285fadd39603ec8a087f30923ca4169b1fc8ed60fbd74579fe53693e357ce79717fcb32a8cb9a5f5dc48bf37b18"

	tcSign  = "0af1958672e6d7a3530588bacf87fa1191d9a075cd38d5651f84a5f4dc97c7605e15e0170fcf27469c29a0c8d6d61b94aeeab01e6caf5d012c3d6ab5eeec2c7d88252fc41db0fd2c27484f37fcb5af8a868b2c0a34771fe667b66a3fd87cca17fe7f0e979861fdfcfa83483bd1a5776a9b10a43c61e91a9a810af30c8fee6391f2005694e109898bbd68f65c7089e934"
	tcHash  = "f74311297112ca6ddfa2602e5c2b8a9d474f944987f71df764d3620b0ee3aa05"
	tcProof = 332
)

func TestMessage(t *testing.T) {
	msgBytes, err := os.ReadFile("message.json")
	if err != nil {
		t.Error(err)
		return
	}

	params := NewSettings(&SSettings{
		FMessageSizeBytes: (2 << 10),
		FWorkSizeBits:     testutils.TCWorkSize,
	})
	msg := LoadMessage(params, msgBytes)
	if msg == nil {
		t.Error("failed load message")
		return
	}

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
