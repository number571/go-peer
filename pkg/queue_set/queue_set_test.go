package queue_set

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/number571/go-peer/pkg/encoding"
)

func TestSettings(t *testing.T) {
	t.Parallel()

	for i := 0; i < 1; i++ {
		testSettings(t, i)
	}
}

func testSettings(t *testing.T, n int) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()
	switch n {
	case 0:
		_ = NewSettings(&SSettings{})
	}
}

func TestQueueSet(t *testing.T) {
	queueSet := NewQueueSet(
		NewSettings(&SSettings{
			FCapacity: 3,
		}),
	)

	sett := queueSet.GetSettings()
	if sett.GetCapacity() != 3 {
		t.Error("got invalid value from settings")
		return
	}

	if _, ok := queueSet.Load([]byte("unknown-key")); ok {
		t.Error("success load unknown key")
		return
	}

	for i := 0; i < 3; i++ {
		key := encoding.Uint64ToBytes(uint64(i))
		if ok := queueSet.Push(key[:], []byte(fmt.Sprintf("_%d_", i))); !ok {
			t.Errorf("failed push %d", i)
			return
		}
	}

	for i := 0; i < 3; i++ {
		key := encoding.Uint64ToBytes(uint64(i))
		val, ok := queueSet.Load(
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

	key1 := encoding.Uint64ToBytes(1)
	if ok := queueSet.Push(key1[:], []byte(fmt.Sprintf("_%d_", 1))); ok {
		t.Error("success push already exist value")
		return
	}

	// start cycle of queue
	i := uint64(4)
	key2 := encoding.Uint64ToBytes(i)
	if ok := queueSet.Push(key2[:], []byte(fmt.Sprintf("_%d_", i))); !ok {
		t.Errorf("failed push %d", i)
		return
	}

	// try load init value
	i = 0
	key3 := encoding.Uint64ToBytes(i)
	if _, ok := queueSet.Load(key3[:]); ok {
		t.Errorf("success load rewrited value %d", i)
		return
	}
}
