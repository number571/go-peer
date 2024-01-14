package utils

import (
	"errors"
	"fmt"
	"testing"
)

var (
	tgErrorn1  = errors.New("error#1")
	tgErrorn2  = errors.New("error#2")
	tgErrorn3  = errors.New("error#3")
	tgErrorn4  = errors.New("error#4")
	tgErrorStr = fmt.Sprintf("%s: %s: %s", tgErrorn1.Error(), tgErrorn2.Error(), tgErrorn3.Error())
)

func TestMergeErrors(t *testing.T) {
	t.Parallel()

	errList := []error{tgErrorn1, nil, tgErrorn2, tgErrorn3}
	err := MergeErrors(errList...)
	if err == nil {
		t.Error("nothing errors?")
		return
	}
	if !errors.Is(err, tgErrorn1) {
		t.Error("not found error#1")
		return
	}
	if !errors.Is(err, tgErrorn2) {
		t.Error("not found error#2")
		return
	}
	if !errors.Is(err, tgErrorn3) {
		t.Error("not found error#3")
		return
	}
	if errors.Is(err, tgErrorn4) {
		t.Error("found not exist error#4")
		return
	}
	if err.Error() != tgErrorStr {
		t.Error("got another string error")
		return
	}
}
