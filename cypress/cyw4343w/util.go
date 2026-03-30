package cyw4343w

import (
	"golang.org/x/exp/constraints"
	"unsafe"
)

func roundUp[T constraints.Integer](x T, y T) T {
	if x%y != 0 {
		return x + y - (x % y)
	}
	return x
}

func toSlice[T any](value *T) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(value)), unsafe.Sizeof(*value))
}

// ntoh16 converts a uint16 from network byte order (big-endian) to host byte order.
func ntoh16(v uint16) uint16 {
	return (v >> 8) | (v << 8)
}

// ntoh32 converts a uint32 from network byte order (big-endian) to host byte order.
func ntoh32(v uint32) uint32 {
	return (v >> 24) | ((v >> 8) & 0xFF00) | ((v << 8) & 0xFF0000) | (v << 24)
}

//sigo:extern memcmp memcmp
func memcmp(dst, src unsafe.Pointer, num uintptr) int
