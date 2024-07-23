package lru

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/number571/go-peer/pkg/encoding"
)

func TestLRUCache(t *testing.T) {
	t.Parallel()

	lruCache := NewLRUCache(3)

	if _, ok := lruCache.Get([]byte("unknown-key")); ok {
		t.Error("success load unknown key")
		return
	}

	for i := 0; i < 3; i++ {
		key := encoding.Uint64ToBytes(uint64(i))
		if ok := lruCache.Set(key[:], []byte(fmt.Sprintf("_%d_", i))); !ok {
			t.Errorf("failed push %d", i)
			return
		}
		if lruCache.GetIndex() != uint64((i+1)%3) {
			t.Error("got invalid index")
			return
		}
	}

	for i := 0; i < 3; i++ {
		key := encoding.Uint64ToBytes(uint64(i))
		val, ok := lruCache.Get(
			key[:],
		)
		if !ok {
			t.Errorf("failed load %d", i)
			return
		}
		if !bytes.Equal(val, []byte(fmt.Sprintf("_%d_", i))) {
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
		if !bytes.Equal(k, key[:]) {
			t.Error("got incorrect key")
			return
		}
	}

	key1 := encoding.Uint64ToBytes(1)
	if ok := lruCache.Set(key1[:], []byte(fmt.Sprintf("_%d_", 1))); ok {
		t.Error("success push already exist value")
		return
	}

	// start cycle of queue
	i := uint64(4)
	key2 := encoding.Uint64ToBytes(i)
	if ok := lruCache.Set(key2[:], []byte(fmt.Sprintf("_%d_", i))); !ok {
		t.Errorf("failed push %d", i)
		return
	}

	// try load init value
	i = 0
	key3 := encoding.Uint64ToBytes(i)
	if _, ok := lruCache.Get(key3[:]); ok {
		t.Errorf("success load rewrited value %d", i)
		return
	}
}
