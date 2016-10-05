package errors

import "fmt"

const (
	// ErrIsClosed is intended to be used when an action is attempted on a closed instance
	ErrIsClosed = Error("cannot perform action on a closed instance")
)

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
