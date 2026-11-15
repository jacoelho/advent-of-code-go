package xmath

import (
	"errors"

	"golang.org/x/exp/constraints"
)

var (
	// ErrNoSolution reports congruences that no integer satisfies together.
	ErrNoSolution = errors.New("xmath: congruences have no common solution")
	// ErrOverflow reports a result that does not fit in the integer type.
	ErrOverflow = errors.New("xmath: result overflows")
)

// CombineCongruences returns the residue and modulus of the integers x with
// x ≡ a (mod m) and x ≡ b (mod n): the modulus is lcm(m, n) and the residue
// lies in [0, lcm(m, n)). The moduli need not be coprime. It returns
// ErrNoSolution when the congruences conflict and ErrOverflow when lcm(m, n)
// does not fit in T. It panics if m or n is not positive.
func CombineCongruences[T constraints.Integer](a, m, b, n T) (residue, modulus T, err error) {
	if m <= 0 || n <= 0 {
		panic("xmath: CombineCongruences moduli must be positive")
	}
	a, b = reduce(a, m), reduce(b, n)
	divisor := GCD(m, n)
	difference := b - a%n
	if b < a%n {
		difference = n - (a%n - b)
	}
	if difference%divisor != 0 {
		return 0, 0, ErrNoSolution
	}
	step := n / divisor
	modulus, ok := CheckedMul(m, step)
	if !ok {
		return 0, 0, ErrOverflow
	}
	inverse, _ := ModInverse(m/divisor, step) // m/divisor and step are coprime.
	// a + k*m with k < step stays below m*step, so it cannot overflow.
	return a + MulMod(difference/divisor, inverse, step)*m, modulus, nil
}
