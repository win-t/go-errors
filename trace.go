package errors

import (
	stderrors "errors"

	"github.com/win-t/go-errors/trace"
)

const traceDeep = 150

type traced struct {
	err   error
	trace []trace.Location
}

func (e *traced) Error() string        { return e.err.Error() }
func (e *traced) Unwrap() error        { return Unwrap(e.err) }
func (e *traced) As(target any) bool   { return As(e.err, target) }
func (e *traced) Is(target error) bool { return Is(e.err, target) }

// Trace will return new error that have stack trace
//
// will return same err if err already have stack trace
//
// use [Is] function to compare the returned error with others, because equality (==) operator will fail
func Trace(err error) error {
	if err == nil {
		return nil
	}

	return doTrace(err)
}

// like [Trace] but suitable for multiple return
func Trace2[A any](a A, err error) (A, error) {
	if err == nil {
		return a, nil
	}

	return a, doTrace(err)
}

// like [Trace] but suitable for multiple return
func Trace3[A, B any](a A, b B, err error) (A, B, error) {
	if err == nil {
		return a, b, nil
	}

	return a, b, doTrace(err)
}

// like [Trace] but suitable for multiple return
func Trace4[A, B, C any](a A, b B, c C, err error) (A, B, C, error) {
	if err == nil {
		return a, b, c, nil
	}

	return a, b, c, doTrace(err)
}

// this function is separated from Trace,
// to make sure that Trace function is simple enough and get inlined
func doTrace(err error) error {
	cur := err
	for cur != nil {
		if _, ok := cur.(*traced); ok {
			return err
		}
		cur = stderrors.Unwrap(cur)
	}

	return &traced{err, trace.Get(1, traceDeep)}
}

// Get stack trace of err
//
// return nil if err doesn't have stack trace
func StackTrace(err error) []trace.Location {
	for err != nil {
		if t, ok := err.(*traced); ok {
			return t.trace
		}
		err = stderrors.Unwrap(err)
	}
	return nil
}
