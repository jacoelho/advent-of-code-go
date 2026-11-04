package matrix

import (
	"cmp"
	"math"
	"math/big"
	"math/bits"
)

// Rat is an immutable arbitrary-precision rational. Its zero value is zero.
// Compare values with Cmp.
type Rat struct {
	_           [0]func()
	numerator   int64
	denominator int64    // Zero denotes 1; otherwise >1 with a nonintegral, possibly unreduced quotient.
	value       *big.Rat // Used only when a reduced numerator or denominator does not fit int64.
}

func integerRat(value int64) Rat { return Rat{numerator: value} }

func NewRat(num, den int64) Rat {
	if den == 1 {
		return integerRat(num)
	}
	return newFraction(num, den)
}

func newFraction(num, den int64) Rat {
	if den == 0 {
		panic("division by zero")
	}
	if num == 0 {
		return Rat{}
	}
	if num == math.MinInt64 && den == -1 {
		return ratFromOwned(big.NewRat(num, den))
	}
	remainder := num % den
	if remainder == 0 {
		return integerRat(num / den)
	}
	numerator, denominator := magnitude(num), magnitude(den)
	divisor := gcdMagnitudes(denominator, magnitude(remainder))
	numerator /= divisor
	denominator /= divisor
	negative := (num < 0) != (den < 0)
	if denominator > math.MaxInt64 || numerator > math.MaxInt64 && !negative {
		return ratFromOwned(big.NewRat(num, den))
	}
	// The negative range includes magnitude 1<<63; converting it already gives MinInt64.
	integer := int64(numerator)
	if negative && integer != math.MinInt64 {
		integer = -integer
	}
	return Rat{numerator: integer, denominator: int64(denominator)}
}

// terminalRat defers reduction for solver results that may only be tested for
// integrality. The determinant is nonzero; integral quotients stay canonical.
func terminalRat(num, den int64) Rat {
	if num == 0 {
		return Rat{}
	}
	if num == math.MinInt64 || den == math.MinInt64 {
		return NewRat(num, den)
	}
	if den < 0 {
		num, den = -num, -den
	}
	if num%den == 0 {
		return integerRat(num / den)
	}
	return Rat{numerator: num, denominator: den}
}

func magnitude(value int64) uint64 {
	if value < 0 {
		return uint64(-(value + 1)) + 1
	}
	return uint64(value)
}

func gcdMagnitudes(a, b uint64) uint64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// NewInteger snapshots value; later mutations cannot change the rational.
func NewInteger(value *big.Int) Rat {
	if value.IsInt64() {
		return integerRat(value.Int64())
	}
	return Rat{value: new(big.Rat).SetInt(value)}
}

// ratFromOwned takes a fresh result; no mutable aliases may survive this call.
func ratFromOwned(value *big.Rat) Rat {
	if value.Num().IsInt64() && value.Denom().IsInt64() {
		result := integerRat(value.Num().Int64())
		if denominator := value.Denom().Int64(); denominator != 1 {
			result.denominator = denominator
		}
		return result
	}
	return Rat{value: value}
}

func (r Rat) rational(scratch *big.Rat) *big.Rat {
	if r.value != nil {
		return r.value
	}
	if r.denominator == 0 {
		return scratch.SetInt64(r.numerator)
	}
	return scratch.SetFrac64(r.numerator, r.denominator)
}

func (r Rat) compactDenominator() int64 {
	if r.denominator == 0 {
		return 1
	}
	return r.denominator
}

// reducedCompact gives canonical observations of a compact value without
// changing an immutable terminal result whose reduction was deferred.
func (r Rat) reducedCompact() (numerator, denominator int64) {
	if r.denominator == 0 {
		return r.numerator, 1
	}
	divisor := int64(gcdMagnitudes(magnitude(r.numerator), uint64(r.denominator)))
	return r.numerator / divisor, r.denominator / divisor
}

// commonNumerators requires two compact values and checks intermediate products.
func (r Rat) commonNumerators(other Rat) (left, right, denominator int64, ok bool) {
	a, b := r.compactDenominator(), other.compactDenominator()
	if a == b {
		return r.numerator, other.numerator, a, true
	}
	divisor := int64(gcdMagnitudes(uint64(a), uint64(b)))
	left, ok = multiplyIntegers(r.numerator, b/divisor)
	if !ok {
		return 0, 0, 0, false
	}
	right, ok = multiplyIntegers(other.numerator, a/divisor)
	if !ok {
		return 0, 0, 0, false
	}
	denominator, ok = multiplyIntegers(a, b/divisor)
	return left, right, denominator, ok
}

func (r Rat) Add(other Rat) Rat {
	if r.value == nil && other.value == nil && r.denominator == 0 && other.denominator == 0 {
		if result, ok := addIntegers(r.numerator, other.numerator); ok {
			return integerRat(result)
		}
	}
	if r.IsZero() {
		return other
	}
	if other.IsZero() {
		return r
	}
	if r.value == nil && other.value == nil && (r.denominator != 0 || other.denominator != 0) {
		if left, right, denominator, ok := r.commonNumerators(other); ok {
			if result, ok := addIntegers(left, right); ok {
				return NewRat(result, denominator)
			}
		}
	}
	var left, right big.Rat
	return ratFromOwned(new(big.Rat).Add(r.rational(&left), other.rational(&right)))
}

func (r Rat) Sub(other Rat) Rat {
	if r.value == nil && other.value == nil && r.denominator == 0 && other.denominator == 0 {
		if result, ok := subtractIntegers(r.numerator, other.numerator); ok {
			return integerRat(result)
		}
	}
	if other.IsZero() {
		return r
	}
	if r.IsZero() {
		return other.Neg()
	}
	if r.value == nil && other.value == nil && (r.denominator != 0 || other.denominator != 0) {
		if left, right, denominator, ok := r.commonNumerators(other); ok {
			if result, ok := subtractIntegers(left, right); ok {
				return NewRat(result, denominator)
			}
		}
	}
	var left, right big.Rat
	return ratFromOwned(new(big.Rat).Sub(r.rational(&left), other.rational(&right)))
}

func (r Rat) Mul(other Rat) Rat {
	if r.value == nil && other.value == nil && r.denominator == 0 && other.denominator == 0 {
		if numerator, ok := multiplyIntegers(r.numerator, other.numerator); ok {
			return integerRat(numerator)
		}
	}
	if r.IsZero() || other.IsZero() {
		return Rat{}
	}
	if r.value == nil && r.denominator == 0 && r.numerator == 1 {
		return other
	}
	if other.value == nil && other.denominator == 0 && other.numerator == 1 {
		return r
	}
	if r.value == nil && other.value == nil && (r.denominator != 0 || other.denominator != 0) {
		if numerator, ok := multiplyIntegers(r.numerator, other.numerator); ok {
			if denominator, ok := multiplyIntegers(r.compactDenominator(), other.compactDenominator()); ok {
				return NewRat(numerator, denominator)
			}
		}
	}
	var left, right big.Rat
	return ratFromOwned(new(big.Rat).Mul(r.rational(&left), other.rational(&right)))
}

func multiplyIntegers(a, b int64) (int64, bool) {
	high, low := bits.Mul64(uint64(a), uint64(b))
	// Correct the unsigned high word for sign extension of each operand.
	high -= uint64(a>>63) & uint64(b)
	high -= uint64(b>>63) & uint64(a)
	result := int64(low)
	return result, high == uint64(result>>63)
}

func subtractIntegers(a, b int64) (int64, bool) {
	result := a - b
	return result, !(b > 0 && result > a || b < 0 && result < a)
}

func addIntegers(a, b int64) (int64, bool) {
	result := a + b
	return result, !(b > 0 && result < a || b < 0 && result > a)
}

func (r Rat) Div(other Rat) Rat {
	if other.IsZero() {
		panic("division by zero")
	}
	if r.value == nil && other.value == nil && r.denominator == 0 && other.denominator == 0 {
		return NewRat(r.numerator, other.numerator)
	}
	if r.IsZero() {
		return Rat{}
	}
	if other.value == nil && other.denominator == 0 && other.numerator == 1 {
		return r
	}
	if r.value == nil && other.value == nil {
		if numerator, ok := multiplyIntegers(r.numerator, other.compactDenominator()); ok {
			if denominator, ok := multiplyIntegers(r.compactDenominator(), other.numerator); ok {
				return NewRat(numerator, denominator)
			}
		}
	}
	var left, right big.Rat
	return ratFromOwned(new(big.Rat).Quo(r.rational(&left), other.rational(&right)))
}

func (r Rat) Neg() Rat {
	if r.value == nil && r.numerator != math.MinInt64 {
		return Rat{numerator: -r.numerator, denominator: r.denominator}
	}
	var scratch big.Rat
	return ratFromOwned(new(big.Rat).Neg(r.rational(&scratch)))
}

func (r Rat) IsZero() bool {
	return r.value == nil && r.numerator == 0
}

func (r Rat) IsInt() bool {
	if r.value == nil {
		return r.denominator == 0
	}
	return r.value.IsInt()
}

func (r Rat) Sign() int {
	if r.value == nil {
		return cmp.Compare(r.numerator, 0)
	}
	return r.value.Sign()
}

func (r Rat) Cmp(other Rat) int {
	if r.value == nil && other.value == nil {
		if r.denominator == other.denominator {
			return cmp.Compare(r.numerator, other.numerator)
		}
		if left, ok := multiplyIntegers(r.numerator, other.compactDenominator()); ok {
			if right, ok := multiplyIntegers(other.numerator, r.compactDenominator()); ok {
				return cmp.Compare(left, right)
			}
		}
	}
	var left, right big.Rat
	return r.rational(&left).Cmp(other.rational(&right))
}

// Numerator returns an independent copy of the reduced numerator.
func (r Rat) Numerator() *big.Int {
	if r.value == nil {
		numerator, _ := r.reducedCompact()
		return big.NewInt(numerator)
	}
	return new(big.Int).Set(r.value.Num())
}

// Denominator returns an independent copy of the positive reduced denominator.
func (r Rat) Denominator() *big.Int {
	if r.value == nil {
		_, denominator := r.reducedCompact()
		return big.NewInt(denominator)
	}
	return new(big.Int).Set(r.value.Denom())
}
