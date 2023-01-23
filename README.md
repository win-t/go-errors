# go-errors.

[![GoReference](https://pkg.go.dev/badge/github.com/win-t/go-errors)](https://pkg.go.dev/github.com/win-t/go-errors)

This package can be used as drop-in replacement for standard errors package.

This package provide `func StackTrace(error) []trace.Location` to get the stack trace.

Stack trace can be attached to any `error` by passing it to `func Trace(error) error`.

`New`, and `Errorf` function will return error that have stack trace.
