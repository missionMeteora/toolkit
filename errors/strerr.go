package errors

import "fmt"

// Error is a constant error type
type Error string

// Error is the error interface implementation
func (err Error) Error() string {
	return string(err)
}

// Fmt returns a formated error based on the current error
func (err Error) Fmt(args ...interface{}) error {
	return Error(fmt.Sprintf(string(err), args...))
}
