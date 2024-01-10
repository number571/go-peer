package echain

import (
	"errors"
	"strings"
)

type link struct {
	err  error
	next *link
}

// ErrorChain will chain multiple errors
type ErrorChain struct {
	head *link
	tail *link
}

// New create a new error chain
func New() *ErrorChain {
	return &ErrorChain{}
}
func (e *ErrorChain) Error() string {
	errs := []string{}
	h := e.head
	for h != nil {
		errs = append(errs, h.err.Error())
		h = h.next
	}
	return strings.Join(errs, ": ")
}

// Unwrap will give the next error
func (e *ErrorChain) Unwrap() error {
	if e.head.next == nil {
		return nil
	}
	ec := &ErrorChain{
		head: e.head.next,
		tail: e.tail,
	}
	return ec
}

// Is will comapre the target
func (e *ErrorChain) Is(target error) bool {
	return errors.Is(e.head.err, target)
}

// Add will place another error in the chain
func (e *ErrorChain) Add(err error) {
	l := &link{
		err: err,
	}
	if e.head == nil {
		e.head = l
		e.tail = l
		return
	}
	e.tail.next = l
	e.tail = l
}

// Errors will return the errors in the chain
func (e *ErrorChain) Errors() []error {
	errs := []error{}
	l := e.head
	for l != nil {
		errs = append(errs, l.err)
		l = l.next
	}
	return errs
}
