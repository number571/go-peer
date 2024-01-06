package testutils

import (
	_ "embed"
	"errors"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/crypto/random"
)

func TestTryN(t *testing.T) {
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
			if random.NewStdPRNG().GetBool() {
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
