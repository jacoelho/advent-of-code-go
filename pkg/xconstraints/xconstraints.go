// Package xconstraints defines type constraints for generic numeric code.
package xconstraints

import "golang.org/x/exp/constraints"

// Number permits any integer, floating-point, or complex type.
type Number interface {
	constraints.Integer | constraints.Float | constraints.Complex
}
