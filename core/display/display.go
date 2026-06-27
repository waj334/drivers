package display

import (
	"unsafe"

	"pkg.si-go.dev/chip/core/hal"
)

const (
	ErrInvalidSurface     hal.Error = "invalid surface"
	ErrInvalidConfig      hal.Error = "invalid config"
	ErrInvalidPixelFormat hal.Error = "invalid pixel format"
)

type PixelFormat uint8

const (
	RGB565 PixelFormat = iota
	RGB888
	XRGB8888
	ARGB8888
)

func BytesPerPixel(format PixelFormat) (int, error) {
	switch format {
	case RGB565:
		return 2, nil
	case RGB888:
		return 3, nil
	case XRGB8888, ARGB8888:
		return 4, nil
	default:
		return 0, ErrInvalidPixelFormat
	}
}

type Rect struct {
	X int
	Y int
	W int
	H int
}

type Surface struct {
	Ptr    unsafe.Pointer
	Len    uintptr
	Width  int
	Height int
	Stride int // bytes per row
	Format PixelFormat
}

func (s Surface) Bytes() []byte {
	return unsafe.Slice((*byte)(s.Ptr), s.Len)
}

type Panel interface {
	Init() error
	Reset() error
	SleepOut() error
	SleepIn() error
	DisplayOn() error
	DisplayOff() error
	SetBacklight(percent uint8) error
}

type ScanoutDisplay interface {
	Panel

	Bounds() Rect
	Surface() Surface
	SetSurface(surface Surface) error
	Present() error
}

type FlushDisplay interface {
	Panel

	Bounds() Rect
	Format() PixelFormat
	Flush(r Rect, pixels []byte, stride int, format PixelFormat) error
}
