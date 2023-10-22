package types

import (
	"fmt"
	"testing"
)

type tsCommand tsCloser
type tsCloser struct {
	fFlag bool
}

func TestStopper(t *testing.T) {
	t.Parallel()

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

func TestCommand(t *testing.T) {
	t.Parallel()

	err := StopAll([]ICommand{
		testNewCommand(false),
		testNewCommand(false),
		testNewCommand(false),
	})
	if err != nil {
		t.Error(err)
		return
	}

	if err := StopAll([]ICommand{testNewCommand(true)}); err == nil {
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

func testNewCommand(flag bool) ICommand {
	return &tsCommand{flag}
}

func (c *tsCommand) Run() error {
	if c.fFlag {
		return fmt.Errorf("some error")
	}
	return nil
}

func (c *tsCommand) Stop() error {
	if c.fFlag {
		return fmt.Errorf("some error")
	}
	return nil
}
