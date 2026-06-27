package primitives

type Rect[T UnitType] struct {
	Point[T]
	Size[T]
}

type RectI32 = Rect[int32]
type RectI64 = Rect[int64]

type RectF32 = Rect[float32]
type RectF64 = Rect[float64]
