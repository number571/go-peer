package types

import (
	"fmt"
	"testing"
)

type tsCloser struct {
	fFlag bool
}

func TestStopper(t *testing.T) {
	err := CloseAll([]ICloser{
		testNewCloser(false),
		testNewCloser(false),
		testNewCloser(false),
	})
	if err != nil {
		t.Error(err)
		return
	}

	if err := CloseAll([]ICloser{testNewCloser(true)}); err == nil {
		t.Error("nothing error?")
		return
	}
}

func testNewCloser(flag bool) ICloser {
	return &tsCloser{flag}
}

func (c *tsCloser) Close() error {
	if c.fFlag {
		return fmt.Errorf("some error")
	}
	return nil
}
