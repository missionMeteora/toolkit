package atomicBool

import "testing"

func TestAtomicBool(t *testing.T) {
	var b Bool
	if b.Get() {
		t.Fatal("invalid state")
	}

	if !b.True() {
		t.Fatal("state change not triggered")
	}

	if !b.Get() {
		t.Fatal("invalid state")
	}

	if b.True() {
		t.Fatal("state change triggered when it shouldn't have")
	}

	if !b.False() {
		t.Fatal("state change not triggered")
	}

	if b.Get() {
		t.Fatal("invalid state")
	}
}
