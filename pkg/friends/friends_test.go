package friends

import (
	"testing"

	"github.com/number571/go-peer/internal/testutils"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

var (
	tgPubKeys = [3]asymmetric.IPubKey{
		asymmetric.LoadRSAPubKey(testutils.TgPubKeys[0]),
		asymmetric.LoadRSAPubKey(testutils.TgPubKeys[1]),
		asymmetric.LoadRSAPubKey(testutils.TgPubKeys[2]),
	}
)

func TestFriends(t *testing.T) {
	f2f := NewF2F()

	f2f.Append(tgPubKeys[0])
	f2f.Append(tgPubKeys[1])
	f2f.Append(tgPubKeys[2])

	if len(f2f.List()) != 3 {
		t.Error("len f2f list != 3")
		return
	}

	f2f.Remove(tgPubKeys[1])
	if len(f2f.List()) != 2 {
		t.Error("len f2f list != 2 after remove")
		return
	}

	if f2f.InList(tgPubKeys[1]) {
		t.Error("deleted value in list")
		return
	}

	if !f2f.InList(tgPubKeys[2]) {
		t.Error("undefined exist value")
		return
	}
}
