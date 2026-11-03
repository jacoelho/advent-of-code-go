package matrix

import "math/big"

// Rat is an immutable arbitrary-precision rational. Its zero value is zero.
// Compare values with Cmp.
type Rat struct {
	value *big.Rat
	_     [0]func()
}

func NewRat(num, den int64) Rat {
	return Rat{value: big.NewRat(num, den)}
}

// NewInteger snapshots value; later mutations cannot change the rational.
func NewInteger(value *big.Int) Rat {
	return Rat{value: new(big.Rat).SetInt(value)}
}

func (r Rat) rational() *big.Rat {
	if r.value == nil {
		return new(big.Rat)
	}
	return r.value
}

func (r Rat) Add(other Rat) Rat {
	return Rat{value: new(big.Rat).Add(r.rational(), other.rational())}
}

func (r Rat) Sub(other Rat) Rat {
	return Rat{value: new(big.Rat).Sub(r.rational(), other.rational())}
}

func (r Rat) Mul(other Rat) Rat {
	return Rat{value: new(big.Rat).Mul(r.rational(), other.rational())}
}

func (r Rat) Div(other Rat) Rat {
	return Rat{value: new(big.Rat).Quo(r.rational(), other.rational())}
}

func (r Rat) Neg() Rat {
	return Rat{value: new(big.Rat).Neg(r.rational())}
}

func (r Rat) IsZero() bool {
	return r.Sign() == 0
}

func (r Rat) IsInt() bool {
	return r.rational().IsInt()
}

func (r Rat) Sign() int {
	return r.rational().Sign()
}

func (r Rat) Cmp(other Rat) int {
	return r.rational().Cmp(other.rational())
}

// Numerator returns an independent copy.
func (r Rat) Numerator() *big.Int {
	return new(big.Int).Set(r.rational().Num())
}

// Denominator returns an independent copy.
func (r Rat) Denominator() *big.Int {
	return new(big.Int).Set(r.rational().Denom())
}
