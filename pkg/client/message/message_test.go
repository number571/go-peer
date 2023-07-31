package message

import (
	"bytes"
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/encoding"
	testutils "github.com/number571/go-peer/test/_data"
)

const (
	tcSalt    = "f2b7a78445fda605506433ea5627ca616af11c1412ce281ac528bb435e6ea8c7fc2426ed7576d0037771e21c7a26cf4e"
	tcSession = "25d4ce9eaa99e1b9fb3ca984ca68218e046c159e586755e2ceca4218f9f55a3a79874687931f890b874a384e004b4941c7ab35b84b9ca4c14f5e2d7ac5e18681ce9c29006a6cb7b21fdaab4bced99b7f6e1201182cef84a9021f4f6ec11c84eca1d5a09fd4a657de4331d1b3fd914f7aaa1a423a853c4fe618b3bb6dcbe10324"
	tcSender  = "fd91c67cd49d6fe1a47620ce6dd148c35993ba147e334d7e0d79d5d19377a9d9846e149da1dee38a7c423111c041dbd24290d2258a5c7a714091200c21e4709a0e7bd7f0ff011b3405f24e1aa9353145801d4e613fd48238372e2fb9fd5db5a673bd59b26ee9c65a78711aafebdec7e2726313a166260ab8718ff097f56f1622e5b6aa5115eeaf90dd1112472eb1a60f2f0efad364f912893a20193f"

	tcSign  = "8280c0e4993b1fc5da05c01c244f95e5fe04318644ed2b68410d1daa2f791dcf1dc66eb968435181d9394a565ad834b341abd8f3909e418412e8f581e03449cec1148bcb58ac36fe2b2d3971fd11ed53163feb4727cc190e992a8a58aded64ba20e344589ad8171d78f4c6e581ff9793b624b038994d2ef8b1de285a80f60b84eec702eb702e441bb20808be69d6d07e"
	tcHash  = "69437ac7bd5533a1aca6741b9e642fc0d45dbe9a528f09b7d89bcfe7a4e42028"
	tcProof = 717
)

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
