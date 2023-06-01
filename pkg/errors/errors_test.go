package errors

import (
	"fmt"
	"testing"
)

const (
	tcErrorResult = "test_3; test_2; test_1;"
)

type tsErrType1 struct{}
type tsErrType2 struct{}
type tsErrType3 struct{}
type tsErrType4 struct{}

func (p *tsErrType1) Error() string {
	return "1"
}

func (p *tsErrType2) Error() string {
	return "2"
}

func (p *tsErrType3) Error() string {
	return "3"
}

func (p *tsErrType4) Error() string {
	return "4"
}

func TestWrapError(t *testing.T) {
	err := WrapError(WrapError(NewError("test_1"), "test_2"), "test_3")
	if err.Error() != tcErrorResult {
		t.Error("result .Error() is invalid")
	}

	if WrapError(nil, "test_1").Error() != NewError("test_1").Error() {
		fmt.Println(WrapError(nil, "test_1").Error())
		fmt.Println(NewError("test_1").Error())
		t.Error("wrap error not equal new error")
	}
}

func TestAppendError(t *testing.T) {
	err1 := AppendError(AppendError(&tsErrType1{}, &tsErrType2{}), &tsErrType3{})
	if !HasError(err1, &tsErrType1{}) || !HasError(err1, &tsErrType2{}) || !HasError(err1, &tsErrType3{}) {
		t.Error("has not error (1, 2 or 3)")
	}
	if HasError(err1, &tsErrType4{}) {
		t.Error("has error (4)")
	}

	if err := AppendError(nil, nil); err != nil {
		t.Error("err2 is not nil")
	}
	if err := AppendError(&tsErrType1{}, nil); !HasError(err, &tsErrType1{}) {
		t.Error("invalid error type")
	}
	if err := AppendError(nil, &tsErrType1{}); !HasError(err, &tsErrType1{}) {
		t.Error("invalid error type")
	}

	err2 := AppendError(nil, &tsErrType1{})
	err3 := AppendError(nil, &tsErrType2{})
	if err := AppendError(&tsErrType3{}, AppendError(err2, err3)); !HasError(err, &tsErrType2{}) {
		t.Error("not found error in error stack")
	}
}
