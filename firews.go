package fivews

import (
	"runtime"
	"strconv"
	"strings"
)

type fivewsError struct {
	message string
	errs    []error
	stack   stack
}

func New(message string) error {
	return &fivewsError{message: message, stack: caller()}
}

func Wrap(message string, cause error) error {
	return &fivewsError{message: message, errs: []error{cause}, stack: caller()}
}

func Join(message string, cause1 error, cause2 error, causes ...error) error {
	errs := make([]error, 0, 2+len(causes))
	errs = append(errs, cause1, cause2)
	errs = append(errs, causes...)
	return &fivewsError{message: message, errs: errs, stack: caller()}
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
	b.WriteString(e.stack.String())
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

type stack []uintptr

func caller() stack {
	var cs [32]uintptr
	n := runtime.Callers(3, cs[:])
	return stack(cs[0:n])
}

type frame struct {
	pc uintptr
}

func (s stack) Frames() []frame {
	frames := make([]frame, len(s))
	for i, pc := range s {
		frames[i] = frame{pc}
	}
	return frames
}

func (s stack) String() string {
	var b strings.Builder
	for _, frame := range s.Frames() {
		b.WriteString(frame.String())
		b.WriteRune('\n')
	}
	return b.String()
}

func (f frame) Function() string {
	fn := runtime.FuncForPC(f.pc)
	if fn == nil {
		return "unknown"
	}
	return fn.Name()
}

func (f frame) FileLine() (string, int) {
	fn := runtime.FuncForPC(f.pc)
	if fn == nil {
		return "unknown", 0
	}
	file, line := fn.FileLine(f.pc)
	return file, line
}

func (f frame) String() string {
	fn := runtime.FuncForPC(f.pc)
	if fn == nil {
		return "unknown"
	}

	file, line := fn.FileLine(f.pc)
	return fn.Name() + " " + file + ":" + strconv.Itoa(line)
}
