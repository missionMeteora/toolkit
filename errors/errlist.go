package errors

import "sync"

// ErrorList is used to chain a list of potential errors
type ErrorList struct {
	mux  sync.RWMutex
	errs []error
}

// Error will return the string-form of the errors
// Note - This is not thread-safe, please run AFTER all pushes are complete
func (e ErrorList) Error() string {
	if len(e.errs) == 0 {
		return ""
	}

	if len(e.errs) == 1 {
		return e.errs[0].Error()
	}

	b := []byte("the following errors occured:\n")
	for _, err := range e.errs {
		b = append(b, err.Error()...)
		b = append(b, '\n')
	}

	return string(b)
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

// Len will return the length of the inner errors list
func (e *ErrorList) Len() (n int) {
	e.mux.RLock()
	if e != nil {
		n = len(e.errs)
	}
	e.mux.RUnlock()
	return
}
