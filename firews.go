package fivews

import (
	"runtime/debug"
	"strings"
)

type fivewsError struct {
	message string
	errs    []error
	stack   *stack
}

// TODO stack traceの取り方見直す
type stack struct {
	stack []byte
}

func New(message string) error {
	return &fivewsError{message: message, stack: &stack{stack: debug.Stack()}}
}

func Wrap(message string, cause error) error {
	return &fivewsError{message: message, errs: []error{cause}, stack: &stack{stack: debug.Stack()}}
}

func Join(message string, cause1 error, cause2 error, causes ...error) error {
	errs := make([]error, 0, 2+len(causes))
	errs = append(errs, cause1, cause2)
	errs = append(errs, causes...)
	return &fivewsError{message: message, errs: errs, stack: &stack{stack: debug.Stack()}}
}

func LastMessage(err error) (bool, string) {
	if err == nil {
		return false, ""
	}
	if fwErr, ok := err.(*fivewsError); ok {
		return true, fwErr.message
	}
	return true, err.Error()
}

// Error returns the error message with stacktrace.
func (e *fivewsError) Error() string {
	var b strings.Builder
	b.WriteString(e.message)
	for _, err := range e.errs {
		b.WriteString(": ")
		if fwErr, ok := err.(*fivewsError); ok {
			// suppress stack trace if the error is a firewsError
			b.WriteString(fwErr.String())
		} else {
			b.WriteString(err.Error())
		}
	}
	b.WriteRune('\n')
	b.Write(e.stack.stack)
	return b.String()
}

// String returns the error message without stacktrace.
func (e *fivewsError) String() string {
	var b strings.Builder
	b.WriteString(e.message)
	for _, err := range e.errs {
		b.WriteString(": ")
		if fwErr, ok := err.(*fivewsError); ok {
			b.WriteString(fwErr.message)
		} else {
			b.WriteString(err.Error())
		}
	}
	return b.String()
}

func (e *fivewsError) Unwrap() []error {
	return e.errs
}
