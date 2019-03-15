package errors

import (
	"strings"
)

// Error represent the wrapped error
type Error interface {
	error

	// Cause return the error that cause this error
	Cause() error

	// StackTrace retrun the stack trace when this error created
	StackTrace() []Location

	// String representation of Error
	String() string

	// internal is just empty function, the purpose is to make this interface cannot be implemented outside this package
	internal()
}
type errorType struct {
	text       string
	cause      error
	stackTrace []Location
}

func (e *errorType) Error() string {
	if e.text != "" {
		return e.text
	}
	if e.cause != nil {
		return e.cause.Error()
	}
	return ""
}

func (e *errorType) Cause() error {
	return e.cause
}

func (e *errorType) StackTrace() []Location {
	return e.stackTrace
}

func (e *errorType) String() string {
	var buff strings.Builder
	var cause error = e
	var first = true

	for cause != nil {
		if first {
			first = false
		} else {
			buff.WriteString("\nCaused by ")
		}
		buff.WriteString("Error: ")
		buff.WriteString(cause.Error())
		buff.WriteString("\n")
		if err, ok := cause.(Error); ok {
			for _, l := range err.StackTrace() {
				buff.WriteString("- ")
				buff.WriteString(l.String())
				buff.WriteString("\n")
			}
			cause = err.Cause()
		} else {
			cause = nil
		}
	}

	return buff.String()
}

func (e *errorType) internal() {}

func new(skip int, text string, err error) Error {
	ret := &errorType{
		text:       text,
		cause:      err,
		stackTrace: generateStackTrace(skip+1, 20),
	}

	return ret
}

func wrap(skip int, text string, err error) Error {
	if err == nil {
		return nil
	}
	if e, ok := err.(Error); ok {
		return e
	}
	return new(skip+1, text, err)
}

// Wrap the err, if err is nil, then return nil
func Wrap(err error) Error {
	return wrap(1, "", err)
}

// New returns an Error that formats as the given text.
func New(text string) Error {
	return new(1, text, nil)
}

// NewWithCause returns an Error that formats as the given text,
// it also indicate that this Error is caused by err.
func NewWithCause(text string, err error) Error {
	return new(1, text, err)
}

// WrapAndCheck do the same thing as `Check(Wrap(err))`
func WrapAndCheck(err error) {
	Check(wrap(1, "", err))
}

// Fail do the same thing as `Check(NewWithCause(text, err))`.
func Fail(text string, err error) {
	Check(new(1, text, err))
}

// CheckOrFail do the same thing as Fail, but only panic when err is not nil
func CheckOrFail(text string, err error) {
	if err != nil {
		Check(new(1, text, err))
	}
}

// Format the error as string
func Format(err error) string {
	if err == nil {
		return ""
	}

	if err2, ok := err.(Error); ok {
		return err2.String()
	}

	return err.Error()
}

// RealCause return the first error (the leaf error) that chain these err
func RealCause(err error) error {
	last := err
	for {
		if err == nil {
			return last
		}

		if err2, ok := err.(Error); ok {
			last = err
			err = err2.Cause()
		} else {
			return err
		}
	}
}
