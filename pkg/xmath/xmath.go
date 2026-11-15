// Package xmath provides integer arithmetic, including overflow-checked and
// modular operations, and geometry helpers.
package xmath

import (
	"golang.org/x/exp/constraints"
)

// Abs returns the absolute value of v. The minimum value of T has no positive
// counterpart and is returned unchanged.
func Abs[T constraints.Signed](v T) T {
	if v < 0 {
		return -v
	}
	return v
}

// Modulo returns a modulo b with the sign of b, so the result lies in [0, b)
// for positive b, unlike the % operator. It panics if b is zero.
func Modulo[T constraints.Signed](a, b T) T {
	return (a%b + b) % b
}

// GCD returns the non-negative greatest common divisor of a and b. GCD(0, 0)
// is 0.
func GCD[T constraints.Integer](a, b T) T {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// LCM returns the least common multiple of value and values. The result is
// not normalised to be non-negative, and it panics if the running result and
// the next value are both zero.
func LCM[T constraints.Integer](value T, values ...T) T {
	result := value
	for _, v := range values {
		result = result * v / GCD(result, v)
	}
	return result
}

// Point is a two-dimensional point with signed integer coordinates.
type Point[T constraints.Signed] interface {
	GetX() T
	GetY() T
}

// PolygonArea returns the absolute shoelace area, truncated to an integer.
// Vertices must be ordered clockwise or counterclockwise. Fewer than three
// vertices have area zero.
func PolygonArea[T constraints.Signed, P Point[T]](vertices []P) T {
	if len(vertices) < 3 {
		return 0
	}

	var area T
	for i := range vertices {
		j := (i + 1) % len(vertices)
		area += vertices[i].GetX() * vertices[j].GetY()
		area -= vertices[j].GetX() * vertices[i].GetY()
	}

	return Abs(area) / 2
}
