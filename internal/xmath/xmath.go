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
