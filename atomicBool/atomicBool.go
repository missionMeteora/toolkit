package atomicBool

import (
	"sync/atomic"
	"unsafe"
)

// Bool is an atomic bool type
type Bool int32

const (
	// False state
	False int32 = iota
	// True state
	True
)

func (b *Bool) getIntPtr() *int32 {
	return (*int32)(unsafe.Pointer(b))
}

// Get will get the current state
func (b *Bool) Get() (state bool) {
	return atomic.LoadInt32(b.getIntPtr()) == True
}

// True will set to true
func (b *Bool) True() (changed bool) {
	return atomic.CompareAndSwapInt32(b.getIntPtr(), False, True)
}

// False will set to false
func (b *Bool) False() (changed bool) {
	return atomic.CompareAndSwapInt32(b.getIntPtr(), True, False)
}
