package layer2

import (
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
)

func TestListCiphersPanic(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()

	cipher1 := symmetric.NewCipherGCM(random.NewRandom().GetBytes(symmetric.CCipherKeySize))
	list := NewKeysContainer(
		symmetric.NewCipherGCM(random.NewRandom().GetBytes(symmetric.CCipherKeySize)),
	)

	_list := list.(*sKeysContainer)
	_list.fMap[cipher1.GetHasher().ToString()] = cipher1

	_ = list.Del(cipher1.GetHasher().ToString())
}

func TestListCiphers(t *testing.T) {
	t.Parallel()

	cipher1 := symmetric.NewCipherGCM(random.NewRandom().GetBytes(symmetric.CCipherKeySize))
	list := NewKeysContainer(
		cipher1,
		symmetric.NewCipherGCM(random.NewRandom().GetBytes(symmetric.CCipherKeySize)),
		symmetric.NewCipherGCM(random.NewRandom().GetBytes(symmetric.CCipherKeySize)),
	)

	if ok := list.Add(cipher1); ok {
		t.Fatal("success add double cipher")
	}
	if ok := list.Del(cipher1.GetHasher().ToString()); !ok {
		t.Fatal("failed delete exist cipher")
	}
	if ok := list.Del(cipher1.GetHasher().ToString()); ok {
		t.Fatal("success delete unknown cipher")
	}

	if len(list.List()) != 2 {
		t.Fatal("invalid length of list")
	}

	for _, v := range list.List() {
		if gotV, ok := list.Get(v.GetHasher().ToString()); !ok || v.ToString() != gotV.ToString() {
			t.Fatal("got invalid value from keys container")
		}
	}
}

func TestContainerTypes(t *testing.T) {
	t.Parallel()

	keysContainer := NewKeysContainer(asymmetric.NewPrivKey().GetPubKey())
	if ok := keysContainer.Add(symmetric.NewCipherCFB(make([]byte, symmetric.CCipherKeySize))); ok {
		t.Fatal("success add another type of key")
	}
	if ok := keysContainer.Add(asymmetric.NewPrivKey().GetPubKey()); !ok {
		t.Fatal("failed add correct type of key")
	}
}
