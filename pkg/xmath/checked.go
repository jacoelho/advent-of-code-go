package xmath

import (
	"math/bits"

	"golang.org/x/exp/constraints"
)

// CheckedAdd returns a+b and whether the sum fits in T. On overflow the
// returned sum is the wrapped value.
func CheckedAdd[T constraints.Integer](a, b T) (T, bool) {
	sum := a + b
	return sum, (b >= 0) == (sum >= a)
}

// CheckedSub returns a-b and whether the difference fits in T. On overflow
// the returned difference is the wrapped value.
func CheckedSub[T constraints.Integer](a, b T) (T, bool) {
	difference := a - b
	return difference, (b >= 0) == (difference <= a)
}

// CheckedMul returns a*b and whether the product fits in T. On overflow the
// returned product is the wrapped value.
func CheckedMul[T constraints.Integer](a, b T) (T, bool) {
	if a == 0 || b == 0 {
		return 0, true
	}
	product := a * b
	// The quotient test alone misses MinInt * -1, whose wrapped product
	// divided by -1 is again MinInt.
	return product, product/b == a && !(a < 0 && b < 0 && product < 0)
}

// MulMod returns a*b modulo m in [0, m) without intermediate overflow.
// Negative operands are reduced modulo m first. It panics if m is not
// positive.
func MulMod[T constraints.Integer](a, b, m T) T {
	if m <= 0 {
		panic("xmath: MulMod modulus must be positive")
	}
	hi, lo := bits.Mul64(uint64(reduce(a, m)), uint64(reduce(b, m)))
	return T(bits.Rem64(hi, lo, uint64(m)))
}

// ModInverse returns x in [0, m) with a*x ≡ 1 (mod m), and false when a and m
// are not coprime. Negative a is reduced modulo m first. It panics if m is not
// positive.
func ModInverse[T constraints.Integer](a, m T) (T, bool) {
	if m <= 0 {
		panic("xmath: ModInverse modulus must be positive")
	}
	// Extended Euclid, with the Bézout coefficient of a kept in [0, m) so
	// that no intermediate value leaves T.
	remainder, nextRemainder := m, reduce(a, m)
	coefficient, nextCoefficient := T(0), T(1)%m
	for nextRemainder != 0 {
		quotient := remainder / nextRemainder
		remainder, nextRemainder = nextRemainder, remainder-quotient*nextRemainder
		step := MulMod(quotient, nextCoefficient, m)
		next := coefficient - step
		if coefficient < step {
			next = m - (step - coefficient)
		}
		coefficient, nextCoefficient = nextCoefficient, next
	}
	if remainder != 1 {
		return 0, false
	}
	return coefficient, true
}

func reduce[T constraints.Integer](a, m T) T {
	a %= m
	if a < 0 {
		a += m
	}
	return a
}
