package primitives

type Size[T UnitType] struct {
	W T
	H T
}

type SizeI32 = Size[int32]
type SizeI64 = Size[int64]

type SizeF32 = Size[float32]
type SizeF64 = Size[float64]
