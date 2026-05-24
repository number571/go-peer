// nolint: err113
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
		t.Fatal("len(r1) != 16 && len(r2) != 16 && len(r3) != 16")
	}

	if !bytes.Equal(r1, r2) || bytes.Equal(r1, r3) {
		t.Fatal("!bytes.Equal(r1, r2) || bytes.Equal(r1, r3)")
	}
}

func TestTryN(t *testing.T) {
	t.Parallel()

	if err := TryN(3, 10*time.Millisecond, func() error { return errors.New("some error") }); err != nil && err.Error() != "some error" { //nolint:err113
		t.Fatal("success tryN with error")
	}
	if err := TryN(3, 10*time.Millisecond, func() error { return nil }); err != nil {
		t.Fatal(err)
	}
	err := TryN(
		1000,
		10*time.Millisecond,
		func() error {
			if random.NewRandom().GetBool() {
				return errors.New("some error") //nolint:err113
			}
			return nil
		},
	)
	if err != nil {
		t.Fatal(err)
	}
}
