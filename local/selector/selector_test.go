package selector

import (
	"testing"

	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/testutils"
)

func TestSelector(t *testing.T) {
	pubKeys := []asymmetric.IPubKey{}
	for _, sPubKey := range testutils.TgPubKeys {
		pubKeys = append(pubKeys, asymmetric.LoadRSAPubKey(sPubKey))
	}

	selector := NewSelector(pubKeys)
	for i := 0; i < 5; i++ {
		checkPubKeys := selector.Shuffle().Return(selector.Length())
		if !testAreUniq(checkPubKeys) {
			t.Error("selector's list's values not unique")
			return
		}
		for i := range pubKeys {
			if pubKeys[i].Address().String() != checkPubKeys[i].Address().String() {
				return
			}
		}
	}

	t.Error("selector's shuffle does not work")
}

func testAreUniq(pubKeys []asymmetric.IPubKey) bool {
	for i := 0; i < len(pubKeys); i++ {
		for j := i + 1; j < len(pubKeys); j++ {
			if pubKeys[i].Address().String() == pubKeys[j].Address().String() {
				return false
			}
		}
	}
	return true
}
