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

func ModInv(a, m int64) int64 {
	a = ((a % m) + m) % m

	if m == 1 {
		return 0
	}

	m0, x0, x1 := m, int64(0), int64(1)

	for a > 1 {
		q := a / m
		m, a = a%m, m
		x0, x1 = x1-q*x0, x0
	}

	if x1 < 0 {
		x1 += m0
	}

	return x1
}

func ModPow(base, exp, mod int64) int64 {
	if mod == 1 {
		return 0
	}

	result := int64(1)
	base = base % mod

	for exp > 0 {
		if exp%2 == 1 {
			result = (result * base) % mod
		}
		exp = exp >> 1
		base = (base * base) % mod
	}

	return result
}
