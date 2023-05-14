package message

import (
	"bytes"
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/encoding"
	testutils "github.com/number571/go-peer/test/_data"
)

const (
	tcSalt    = "1e0803726723f2ec609a048aa7779f772ebf7002489d2e2f7bc6ae91e8eb619aba74e794c2d0a00e16bda23d75076513"
	tcSession = "b53ac6435763f917e6da7cce7d7d37dc658de791f2da58655f1a715ac44d0b55c22342fa1287b6bc0bfadd4ab7ff480dbdfd09d62e861073ffb9291a6b8835499ee290792015faba1ef21c51facbe09c0c0c7454eb132a4219b6781a5673af864649cb8f8054b0744027fb5ba28894864935a45e133677fe49e8e92de42963f4"
	tcSender  = "3ec534f2fa1df35376c55256bb1ee02e8a808d1813a9f4cac7537c232a2cb481987cfab744a65b75e8c707674e9996ef354a5616dd8a194337c2647815a0cbd92514529d32719a7e2b97ee52fe390e7653c88ff652ec53cddd2d6e3aac40c2adffbaf4758b406cfc3ee943c780ada11a2dc602c9fbe02821ce9428150c81f283c42a7f30eeb5ccee05375a5a93c0cbc10c8a3edac448385e0632c5dc"

	tcHash  = "cbf42e03ea0760e86248ce59e822b4abf47239e86cd24d371dd87adf64232bac"
	tcSign  = "2c2ebfe79d3911b9e987b1e27821dc834fb36352e3ac11dc2f590a6163d8f054dd308f84607337820ea3fdca4587c56f83b80d9578d2750ecacb495aae994a37c1ba66d17fcb4b87898ba4e256ffcdc629e68c459c48c5963f61b180ae47d4d44647df08a1ed2ea9458ff2f5400a7ec8121be37835d5e535f61517a987aba8b1fe52703aa7267e239e2a726a2dbf4935"
	tcProof = 1651
)

func TestMessage(t *testing.T) {
	msgBytes, err := os.ReadFile("message.json")
	if err != nil {
		t.Error(err)
		return
	}

	params := NewSettings(&SSettings{
		FWorkSize:    testutils.TCWorkSize,
		FMessageSize: testutils.TCMessageSize,
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

	if !bytes.Equal(msg.GetBody().GetHash(), encoding.HexDecode(tcHash)) {
		t.Error("incorrect hash value")
		return
	}

	if !bytes.Equal(msg.GetBody().GetSign(), encoding.HexDecode(tcSign)) {
		t.Error("incorrect sign value")
		return
	}

	if msg.GetBody().GetProof() != tcProof {
		t.Error("incorrect proof value")
		return
	}

	// pass payload -> message is encrypted
}
