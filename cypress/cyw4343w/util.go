package cyw4343w

import (
	"golang.org/x/exp/constraints"
	"unsafe"
)

func roundUp[T constraints.Integer](x, y T) T {
	if x%y != 0 {
		return x + y - (x % y)
	}
	return x
}

func toSlice[T any](value *T) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(value)), unsafe.Sizeof(*value))
}
