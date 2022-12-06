package closer

import (
	"fmt"
	"testing"

	"github.com/number571/go-peer/modules"
)

type tsCloser struct {
	fFlag bool
}

func TestCloser(t *testing.T) {
	err := CloseAll([]modules.ICloser{
		testNewCloser(false),
		testNewCloser(false),
		testNewCloser(false),
	})
	if err != nil {
		t.Error(err)
		return
	}

	if err := CloseAll([]modules.ICloser{testNewCloser(true)}); err == nil {
		t.Errorf("nothing error?")
		return
	}
}

func testNewCloser(flag bool) modules.ICloser {
	return &tsCloser{flag}
}

func (c *tsCloser) Close() error {
	if c.fFlag {
		return fmt.Errorf("some error")
	}
	return nil
}
