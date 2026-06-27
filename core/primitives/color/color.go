// Package color provides a small, allocation-free color primitive with
// conversions to/from the packed pixel formats used by the STM32 LTDC and
// DMA2D peripherals (ARGB8888, RGB888, RGB565, ARGB1555, ARGB4444).
//
// Color is a value type (4 bytes). All operations take/return by value, never
// allocate, and use no defer, so they are safe to call from any sigo context
// including ISRs. Alpha is straight (non-premultiplied) unless noted otherwise.
package color

import "golang.org/x/exp/constraints"

// Color is an 8-bit-per-channel RGBA color with straight alpha.
type Color struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

// ---------------------------------------------------------------------------
// Constructors
// ---------------------------------------------------------------------------

// RGBA constructs a Color from explicit channel bytes.
func RGBA(r, g, b, a uint8) Color { return Color{r, g, b, a} }

// RGB constructs an opaque Color (A = 255).
func RGB(r, g, b uint8) Color { return Color{r, g, b, 0xFF} }

// Gray constructs an opaque gray color.
func Gray(v uint8) Color { return Color{v, v, v, 0xFF} }

// ColorF constructs a Color from float channels in [0,1]. Inputs are clamped
// and rounded, so values slightly outside the range (e.g. 1.0000001) no longer
// wrap to near-black. For a raw, branch-free fast path use a struct literal.
func ColorF[T constraints.Float](r, g, b, a T) Color {
	return Color{
		R: uint8(clamp01(r)*255 + 0.5),
		G: uint8(clamp01(g)*255 + 0.5),
		B: uint8(clamp01(b)*255 + 0.5),
		A: uint8(clamp01(a)*255 + 0.5),
	}
}

func clamp01[T constraints.Float](v T) T {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

// FromARGB unpacks a 0xAARRGGBB word.
func FromARGB(v uint32) Color {
	return Color{uint8(v >> 16), uint8(v >> 8), uint8(v), uint8(v >> 24)}
}

// FromRGBA unpacks a 0xRRGGBBAA word.
func FromRGBA(v uint32) Color {
	return Color{uint8(v >> 24), uint8(v >> 16), uint8(v >> 8), uint8(v)}
}

// FromRGB unpacks a 0x00RRGGBB word as an opaque color.
func FromRGB(v uint32) Color {
	return Color{uint8(v >> 16), uint8(v >> 8), uint8(v), 0xFF}
}

// FromRGB565 unpacks a 16-bit RGB565 value, expanding each channel to 8 bits
// via bit replication so full-scale values round-trip exactly.
func FromRGB565(v uint16) Color {
	r5 := uint8(v>>11) & 0x1F
	g6 := uint8(v>>5) & 0x3F
	b5 := uint8(v) & 0x1F
	return Color{
		R: r5<<3 | r5>>2,
		G: g6<<2 | g6>>4,
		B: b5<<3 | b5>>2,
		A: 0xFF,
	}
}

// FromARGB1555 unpacks a 16-bit ARGB1555 value.
func FromARGB1555(v uint16) Color {
	r5 := uint8(v>>10) & 0x1F
	g5 := uint8(v>>5) & 0x1F
	b5 := uint8(v) & 0x1F
	a := uint8(0)
	if v&0x8000 != 0 {
		a = 0xFF
	}
	return Color{
		R: r5<<3 | r5>>2,
		G: g5<<3 | g5>>2,
		B: b5<<3 | b5>>2,
		A: a,
	}
}

// FromARGB4444 unpacks a 16-bit ARGB4444 value (4→8 bit replication).
func FromARGB4444(v uint16) Color {
	a4 := uint8(v>>12) & 0xF
	r4 := uint8(v>>8) & 0xF
	g4 := uint8(v>>4) & 0xF
	b4 := uint8(v) & 0xF
	return Color{
		R: r4<<4 | r4,
		G: g4<<4 | g4,
		B: b4<<4 | b4,
		A: a4<<4 | a4,
	}
}

// ---------------------------------------------------------------------------
// Packers — match LTDC / DMA2D pixel formats. Channel downsampling truncates
// the low bits (the hardware behavior), so endpoints round-trip exactly.
// ---------------------------------------------------------------------------

// ARGB packs to 0xAARRGGBB (LTDC/DMA2D ARGB8888).
func (c Color) ARGB() uint32 {
	return uint32(c.A)<<24 | uint32(c.R)<<16 | uint32(c.G)<<8 | uint32(c.B)
}

// RGBA32 packs to 0xRRGGBBAA.
func (c Color) RGBA32() uint32 {
	return uint32(c.R)<<24 | uint32(c.G)<<16 | uint32(c.B)<<8 | uint32(c.A)
}

// RGB packs to 0x00RRGGBB (LTDC/DMA2D RGB888).
func (c Color) RGB() uint32 {
	return uint32(c.R)<<16 | uint32(c.G)<<8 | uint32(c.B)
}

// BGR packs to 0x00BBGGRR, for panels/framebuffers wired in BGR order.
func (c Color) BGR() uint32 {
	return uint32(c.B)<<16 | uint32(c.G)<<8 | uint32(c.R)
}

// RGB565 packs to 16-bit RGB565 (LTDC/DMA2D RGB565).
func (c Color) RGB565() uint16 {
	return uint16(c.R>>3)<<11 | uint16(c.G>>2)<<5 | uint16(c.B>>3)
}

// ARGB1555 packs to 16-bit ARGB1555. Alpha collapses to 1 bit at the >=128
// threshold.
func (c Color) ARGB1555() uint16 {
	var a uint16
	if c.A >= 128 {
		a = 1
	}
	return a<<15 | uint16(c.R>>3)<<10 | uint16(c.G>>3)<<5 | uint16(c.B>>3)
}

// ARGB4444 packs to 16-bit ARGB4444.
func (c Color) ARGB4444() uint16 {
	return uint16(c.A>>4)<<12 | uint16(c.R>>4)<<8 | uint16(c.G>>4)<<4 | uint16(c.B>>4)
}

// ---------------------------------------------------------------------------
// Operations
// ---------------------------------------------------------------------------

// WithAlpha returns a copy with alpha replaced.
func (c Color) WithAlpha(a uint8) Color { c.A = a; return c }

func (c Color) IsOpaque() bool      { return c.A == 0xFF }
func (c Color) IsTransparent() bool { return c.A == 0 }

// Premultiply returns the premultiplied-alpha form of c.
func (c Color) Premultiply() Color {
	a := uint32(c.A)
	return Color{
		R: uint8(uint32(c.R) * a / 255),
		G: uint8(uint32(c.G) * a / 255),
		B: uint8(uint32(c.B) * a / 255),
		A: c.A,
	}
}

// Over composites c onto dst using straight-alpha source-over. The common
// opaque-dst case (a framebuffer pixel) reduces to a 255-scaled blend with no
// per-channel divide beyond the fixed /255.
func (c Color) Over(dst Color) Color {
	sa := uint32(c.A)
	switch sa {
	case 0xFF:
		return c
	case 0:
		return dst
	}
	ia := 255 - sa
	da := uint32(dst.A) * ia / 255
	outA := sa + da
	if outA == 0 {
		return Color{}
	}
	return Color{
		R: uint8((uint32(c.R)*sa + uint32(dst.R)*da) / outA),
		G: uint8((uint32(c.G)*sa + uint32(dst.G)*da) / outA),
		B: uint8((uint32(c.B)*sa + uint32(dst.B)*da) / outA),
		A: uint8(outA),
	}
}

// Lerp linearly interpolates every channel toward b. t is clamped to [0,1].
func (c Color) Lerp(b Color, t float32) Color {
	t = clamp01(t)
	return Color{
		R: lerp8(c.R, b.R, t),
		G: lerp8(c.G, b.G, t),
		B: lerp8(c.B, b.B, t),
		A: lerp8(c.A, b.A, t),
	}
}

func lerp8(a, b uint8, t float32) uint8 {
	return uint8(float32(a) + (float32(b)-float32(a))*t + 0.5)
}

// Scale multiplies RGB brightness by s (clamped), leaving alpha untouched.
func (c Color) Scale(s float32) Color {
	return Color{
		R: scale8(c.R, s),
		G: scale8(c.G, s),
		B: scale8(c.B, s),
		A: c.A,
	}
}

func scale8(v uint8, s float32) uint8 {
	f := float32(v) * s
	if f <= 0 {
		return 0
	}
	if f >= 255 {
		return 255
	}
	return uint8(f + 0.5)
}

// Luminance returns Rec.601 perceptual brightness (weights sum to 256).
func (c Color) Luminance() uint8 {
	return uint8((77*uint32(c.R) + 150*uint32(c.G) + 29*uint32(c.B)) >> 8)
}

// Grayscale collapses to luminance, preserving alpha.
func (c Color) Grayscale() Color {
	y := c.Luminance()
	return Color{y, y, y, c.A}
}

// ---------------------------------------------------------------------------
// Named colors. These are read-only; the ImmutableGlobals pass keeps them in
// flash/.rodata rather than SRAM.
// ---------------------------------------------------------------------------

var (
	Transparent = Color{0, 0, 0, 0}
	Black       = Color{0, 0, 0, 0xFF}
	White       = Color{0xFF, 0xFF, 0xFF, 0xFF}
	Red         = Color{0xFF, 0, 0, 0xFF}
	Green       = Color{0, 0xFF, 0, 0xFF}
	Blue        = Color{0, 0, 0xFF, 0xFF}
	Yellow      = Color{0xFF, 0xFF, 0, 0xFF}
	Cyan        = Color{0, 0xFF, 0xFF, 0xFF}
	Magenta     = Color{0xFF, 0, 0xFF, 0xFF}
)
