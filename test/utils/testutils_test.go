// nolint: goerr113
package testutils

import (
	"bytes"
	_ "embed"
	"errors"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/crypto/random"
)

func TestPseudoRandomBytes(t *testing.T) {
	t.Parallel()

	r1 := PseudoRandomBytes(1)
	r2 := PseudoRandomBytes(1)
	r3 := PseudoRandomBytes(2)

	if len(r1) != 16 && len(r2) != 16 && len(r3) != 16 {
		t.Error("len(r1) != 16 && len(r2) != 16 && len(r3) != 16")
		return
	}

	if !bytes.Equal(r1, r2) || bytes.Equal(r1, r3) {
		t.Error("!bytes.Equal(r1, r2) || bytes.Equal(r1, r3)")
		return
	}
}

func TestTryN(t *testing.T) {
	t.Parallel()

	if err := TryN(3, 10*time.Millisecond, func() error { return errors.New("some error") }); err != nil && err.Error() != "some error" {
		t.Error("success tryN with error")
		return
	}
	if err := TryN(3, 10*time.Millisecond, func() error { return nil }); err != nil {
		t.Error(err)
		return
	}
	err := TryN(
		1000,
		10*time.Millisecond,
		func() error {
			if random.NewCSPRNG().GetBool() {
				return errors.New("some error")
			}
			return nil
		},
	)
	if err != nil {
		t.Error(err)
		return
	}
}
