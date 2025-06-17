package cyw4343w

import "unsafe"

type bufferHeader struct {
	next unsafe.Pointer
	_    [2]byte
}
type commonBusHeader struct {
	bufferHeader
	busHeader [busHeaderLen]byte
}

type cdcHeader struct {
	cmd    uint32 /* ioctl command value */
	len    uint32 /* lower 16: output buflen; upper 16: input buflen (excludes header) */
	flags  uint32 /* flag defns given in bcmcdc.h */
	status uint32 /* status code returned from the device */
}

type controlHeader struct {
	commonBusHeader
	cdcHeader
}
