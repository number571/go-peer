package utils

import (
	"errors"
	"fmt"
	"testing"
)

var (
	err1   = errors.New("error#1")
	err2   = errors.New("error#2")
	err3   = errors.New("error#3")
	err4   = errors.New("error#4")
	errStr = fmt.Sprintf("%s: %s: %s", err1.Error(), err2.Error(), err3.Error())
)

func TestMergeErrors(t *testing.T) {
	t.Parallel()

	errList := []error{err1, nil, err2, err3}
	err := MergeErrors(errList...)
	if err == nil {
		t.Error("nothing errors?")
		return
	}
	if !errors.Is(err, err1) {
		t.Error("not found error#1")
		return
	}
	if !errors.Is(err, err2) {
		t.Error("not found error#2")
		return
	}
	if !errors.Is(err, err3) {
		t.Error("not found error#3")
		return
	}
	if errors.Is(err, err4) {
		t.Error("found not exist error#4")
		return
	}
	if err.Error() != errStr {
		t.Error("got another string error")
		return
	}
}
