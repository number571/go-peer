package symmetric

import (
	"testing"

	"github.com/number571/go-peer/pkg/crypto/random"
)

func TestListCiphersPanic(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()

	cipher1 := NewCipherGCM(random.NewRandom().GetBytes(CCipherKeySize))
	list := NewListCiphers(
		NewCipherGCM(random.NewRandom().GetBytes(CCipherKeySize)),
	)

	_list := list.(*sListCiphers)
	_list.fMap[cipher1.ToString()] = struct{}{}

	_ = list.Del(cipher1)
}

func TestListCiphers(t *testing.T) {
	t.Parallel()

	cipher1 := NewCipherGCM(random.NewRandom().GetBytes(CCipherKeySize))
	list := NewListCiphers(
		cipher1,
		NewCipherGCM(random.NewRandom().GetBytes(CCipherKeySize)),
		NewCipherGCM(random.NewRandom().GetBytes(CCipherKeySize)),
	)

	if ok := list.Add(cipher1); ok {
		t.Fatal("success add double cipher")
	}
	if ok := list.Del(cipher1); !ok {
		t.Fatal("failed delete exist cipher")
	}
	if ok := list.Del(cipher1); ok {
		t.Fatal("success delete unknown cipher")
	}

	if len(list.Get()) != 2 {
		t.Fatal("invalid length of list")
	}
}
