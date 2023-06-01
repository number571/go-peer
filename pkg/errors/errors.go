package errors

import (
	std_errors "errors"
	"fmt"
	"strings"
)

var (
	_ error = &sError{}
)

type sError struct {
	fErrStack []error
}

func (p *sError) Error() string {
	s := ""
	for i := len(p.fErrStack) - 1; i >= 0; i-- {
		s += fmt.Sprintf("%s; ", p.fErrStack[i].Error())
	}
	return strings.TrimSpace(s)
}

func NewError(pMsg string) error {
	return WrapError(nil, pMsg)
}

func WrapError(pErr error, pMsg string) error {
	return AppendError(pErr, std_errors.New(pMsg))
}

func AppendError(pErr, pErr2 error) error {
	if pErr == nil && pErr2 == nil {
		return nil
	}
	if pErr == nil {
		if _, ok := pErr2.(*sError); ok {
			return pErr2
		}
		return &sError{[]error{pErr2}}
	}
	if pErr2 == nil {
		if _, ok := pErr.(*sError); ok {
			return pErr
		}
		return &sError{[]error{pErr}}
	}
	v, ok := pErr.(*sError)
	if ok {
		return &sError{copyAndAppend(v.fErrStack, pErr2)}
	}
	return &sError{[]error{pErr, pErr2}}
}

func HasError(pErr, pErrType error) bool {
	if pErr == nil {
		return false
	}
	v, ok := pErr.(*sError)
	if !ok {
		return std_errors.Is(pErr, pErrType)
	}
	for i := len(v.fErrStack) - 1; i >= 0; i-- {
		vErr, ok := v.fErrStack[i].(*sError)
		if ok {
			if exist := HasError(vErr, pErrType); exist {
				return true
			}
			continue
		}
		if std_errors.Is(v.fErrStack[i], pErrType) {
			return true
		}
	}
	return false
}

func copyAndAppend(slice []error, val error) []error {
	newSlice := make([]error, len(slice))
	copy(newSlice[:], slice)
	return append(newSlice, val)
}
