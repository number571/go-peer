package interrupt

import (
	"fmt"
	"testing"

	"github.com/number571/go-peer/pkg/types"
)

type tsCloser struct {
	fFlag bool
}

func TestCloser(t *testing.T) {
	t.Parallel()

	err := CloseAll([]types.ICloser{
		testNewCloser(false),
		testNewCloser(false),
		testNewCloser(false),
	})
	if err != nil {
		t.Error(err)
		return
	}

	if err := CloseAll([]types.ICloser{testNewCloser(true)}); err == nil {
		t.Error("nothing error?")
		return
	}
}

func testNewCloser(flag bool) types.ICloser {
	return &tsCloser{flag}
}

func (c *tsCloser) Close() error {
	if c.fFlag {
		return fmt.Errorf("some error")
	}
	return nil
}
