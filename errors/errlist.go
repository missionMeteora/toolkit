package errors

// ErrorList is used to chain a list of potential errors
type ErrorList []error

// Error is the error interface implementation
func (e ErrorList) Error() string {
	if len(e) == 0 {
		return ""
	}

	b := []byte("the following errors occured:\n")
	for _, err := range e {
		b = append(b, err.Error()...)
		b = append(b, '\n')
	}

	return string(b)
}

// Append appends an error to the error list
func (e ErrorList) Append(err error) ErrorList {
	if err == nil {
		return e
	}

	if oe, ok := err.(ErrorList); ok {
		return append(e, oe...)
	}
	return append(e, err)
}

// Push adds the error to the list if it is not nil
func (e *ErrorList) Push(err error) {
	if e == nil || err == nil {
		return
	}
	switch err := err.(type) {
	case ErrorList:
		*e = append(*e, err...)
	case *ErrorList:
		*e = append(*e, *err...)
	default:
		*e = append(*e, err)
	}
}

// Err returns error value of ErrorList
func (e ErrorList) Err() error {
	if len(e) == 0 {
		return nil
	}
	return e
}
