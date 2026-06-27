package primitives

type Point[T UnitType] struct {
	X T
	Y T
}

type PointI32 = Point[int32]
type PointI64 = Point[int64]

type PointF32 = Point[float32]
type PointF64 = Point[float64]
