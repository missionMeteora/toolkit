package errors

import (
	"testing"
)

func TestErrorList(t *testing.T) {
	var errs ErrorList
	errs.Push(Error("hello world"))
	if errs.Err() == nil {
		t.Fatal("error is nil when it should not be")
	}

	if errs.Len() != 1 {
		t.Fatal("invalid errorlist length")
	}

	errs.ForEach(func(err error) {
		if err.Error() != "hello world" {
			t.Fatal("invalid error")
		}
	})

	testErrorlistInterface(errs)
}

func testErrorlistInterface(err error) {
	return
}
