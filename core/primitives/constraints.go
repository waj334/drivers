package primitives

import "golang.org/x/exp/constraints"

type UnitType interface {
	constraints.Integer | constraints.Float
}
