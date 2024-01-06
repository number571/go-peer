package asymmetric

import (
	"testing"

	testutils "github.com/number571/go-peer/test/utils"
)

var (
	tgPubKeys = [3]IPubKey{
		LoadRSAPubKey(testutils.TgPubKeys[0]),
		LoadRSAPubKey(testutils.TgPubKeys[1]),
		LoadRSAPubKey(testutils.TgPubKeys[2]),
	}
)

func TestFriends(t *testing.T) {
	t.Parallel()

	list := NewListPubKeys()

	list.AddPubKey(tgPubKeys[0])
	list.AddPubKey(tgPubKeys[1])
	list.AddPubKey(tgPubKeys[2])

	if len(list.GetPubKeys()) != 3 {
		t.Error("len f2f list != 3")
		return
	}

	list.DelPubKey(tgPubKeys[1])
	if len(list.GetPubKeys()) != 2 {
		t.Error("len f2f list != 2 after remove")
		return
	}

	if list.InPubKeys(tgPubKeys[1]) {
		t.Error("deleted value in list")
		return
	}

	if !list.InPubKeys(tgPubKeys[2]) {
		t.Error("undefined exist value")
		return
	}
}
