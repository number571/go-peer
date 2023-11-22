package utils

import (
	"errors"
	"fmt"
	"testing"
)

var (
	tgErrorn1  = errors.New("error#1")
	tgErrorn2  = errors.New("error#2")
	tgErrorStr = fmt.Sprintf("%s, %s", tgErrorn2.Error(), tgErrorn1.Error())
)

func TestMergeErrors(t *testing.T) {
	errList := []error{tgErrorn1, nil, tgErrorn2}
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
	if err.Error() != tgErrorStr {
		fmt.Println(err.Error())
		t.Error("got another string error")
		return
	}
}
