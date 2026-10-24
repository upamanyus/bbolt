package common

import (
	"unsafe"
)

// The helpers below use unsafe.Add and unsafe.Slice rather than uintptr
// arithmetic, so that every pointer they return is derived from base without a
// conversion from uintptr, which a formal model of Go cannot follow. They
// behave as before.

func UnsafeAdd(base unsafe.Pointer, offset uintptr) unsafe.Pointer {
	return unsafe.Add(base, offset)
}

func UnsafeIndex(base unsafe.Pointer, offset uintptr, elemsz uintptr, n int) unsafe.Pointer {
	return unsafe.Add(base, offset+uintptr(n)*elemsz)
}

// UnsafeByteSlice returns the bytes [offset+i, offset+j) from base, as a slice
// of capacity j-i. It slices after unsafe.Slice rather than starting at
// offset+i, so that for i == j at the very end of an allocation it does not
// form a pointer one past the end.
func UnsafeByteSlice(base unsafe.Pointer, offset uintptr, i, j int) []byte {
	return unsafe.Slice((*byte)(unsafe.Add(base, offset)), j)[i:j:j]
}
