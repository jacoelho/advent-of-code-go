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

func Modulo[T constraints.Signed](a, b T) T {
	return (a%b + b) % b
}

func GCD[T constraints.Signed](a, b T) T {
	a = Abs(a)
	b = Abs(b)
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func LCM[T constraints.Signed](value T, values ...T) T {
	result := value
	for _, v := range values {
		result = result * v / GCD(result, v)
	}
	return result
}

type Point[T constraints.Signed] interface {
	GetX() T
	GetY() T
}

// PolygonArea calculates the area of a polygon using the Shoelace formula.
// The polygon vertices should be provided in order (clockwise or counter-clockwise).
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
