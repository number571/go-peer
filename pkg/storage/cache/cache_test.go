package cache

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/number571/go-peer/pkg/encoding"
)

func TestLRUCache(t *testing.T) {
	t.Parallel()

	lruCache := NewLRUCache(3)

	if _, ok := lruCache.Get("unknown-key"); ok {
		t.Fatal("success load unknown key")
	}

	for i := 0; i < 3; i++ {
		key := encoding.Uint64ToBytes(uint64(i)) //nolint:gosec
		if ok := lruCache.Set(string(key[:]), []byte(fmt.Sprintf("_%d_", i))); !ok {
			t.Errorf("failed push %d", i)
			return
		}
		if lruCache.GetIndex() != uint64((i+1)%3) { //nolint:gosec
			t.Fatal("got invalid index")
		}
	}

	for i := 0; i < 3; i++ {
		key := encoding.Uint64ToBytes(uint64(i)) //nolint:gosec
		val, ok := lruCache.Get(
			string(key[:]),
		)
		if !ok {
			t.Errorf("failed load %d", i)
			return
		}
		if !bytes.Equal(val.([]byte), []byte(fmt.Sprintf("_%d_", i))) {
			t.Errorf("value is incorrect %d", i)
			return
		}
	}

	for i := uint64(0); ; i++ {
		k, ok := lruCache.GetKey(i)
		if !ok {
			break
		}
		key := encoding.Uint64ToBytes(i)
		if !bytes.Equal([]byte(k), key[:]) {
			t.Fatal("got incorrect key")
		}
	}

	key1 := encoding.Uint64ToBytes(1)
	if ok := lruCache.Set(string(key1[:]), []byte(fmt.Sprintf("_%d_", 1))); ok {
		t.Fatal("success push already exist value")
	}

	// start cycle of queue
	i := uint64(4)
	key2 := encoding.Uint64ToBytes(i)
	if ok := lruCache.Set(string(key2[:]), []byte(fmt.Sprintf("_%d_", i))); !ok {
		t.Fatalf("failed push %d", i)
	}

	// try load init value
	i = 0
	key3 := encoding.Uint64ToBytes(i)
	if _, ok := lruCache.Get(string(key3[:])); ok {
		t.Fatalf("success load rewrited value %d", i)
	}
}
