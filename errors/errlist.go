package errors

import "sync"

// ErrorList is used to chain a list of potential errors
type ErrorList struct {
	mux  sync.RWMutex
	errs []error
}

// Error will return the string-form of the errors
func (e *ErrorList) Error() (str string) {
	var b []byte
	e.mux.RLock()
	if e == nil || len(e.errs) == 0 {
		goto END
	}

	b = []byte("the following errors occured:\n")
	for _, err := range e.errs {
		b = append(b, err.Error()...)
		b = append(b, '\n')
	}

	str = string(b)

END:
	e.mux.RUnlock()
	return
}

// Err will return an error if the errorlist is not empty
// If the errorlist is empty - nil is returned
func (e *ErrorList) Err() (err error) {
	e.mux.RLock()
	if e != nil && len(e.errs) > 0 {
		err = e
	}
	e.mux.RLock()
	return
}

// Push will push an error to the errorlist
// If the errorlist is nil, it will be created
func (e *ErrorList) Push(err error) {
	if err == nil {
		return
	}

	e.mux.Lock()
	if e == nil {
		*e = ErrorList{}
	}

	e.errs = append(e.errs, err)
	e.mux.Unlock()
}
