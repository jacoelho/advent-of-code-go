package xmath

import (
	"golang.org/x/exp/constraints"
)

func Abs[T constraints.Signed](v T) T {
	if v < 0 {
		return -v
	}
	return v
}
