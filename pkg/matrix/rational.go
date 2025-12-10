package matrix

import "github.com/jacoelho/advent-of-code-go/pkg/xmath"

// Rat represents a rational number using integers
type Rat struct {
	Num, Den int64
}

// NewRat creates a new rational number, normalizing it to reduced form
func NewRat(num, den int64) Rat {
	if den == 0 {
		panic("rational number with zero denominator")
	}
	if den < 0 {
		num, den = -num, -den
	}
	g := xmath.GCD(num, den)
	if g != 0 {
		num /= g
		den /= g
	}
	return Rat{Num: num, Den: den}
}

// Field[Rat] interface implementation
func (r Rat) Add(other Rat) Rat {
	lcm := xmath.LCM(r.Den, other.Den)
	num := r.Num*(lcm/r.Den) + other.Num*(lcm/other.Den)
	return NewRat(num, lcm)
}

func (r Rat) Sub(other Rat) Rat {
	lcm := xmath.LCM(r.Den, other.Den)
	num := r.Num*(lcm/r.Den) - other.Num*(lcm/other.Den)
	return NewRat(num, lcm)
}

func (r Rat) Mul(other Rat) Rat {
	return NewRat(r.Num*other.Num, r.Den*other.Den)
}

func (r Rat) Div(other Rat) Rat {
	if other.Num == 0 {
		panic("division by zero")
	}
	return NewRat(r.Num*other.Den, r.Den*other.Num)
}

func (r Rat) Neg() Rat {
	return Rat{Num: -r.Num, Den: r.Den}
}

func (r Rat) IsZero() bool {
	return r.Num == 0
}

func (r Rat) Zero() Rat {
	return NewRat(0, 1)
}

func (r Rat) One() Rat {
	return NewRat(1, 1)
}

// Numerator returns the numerator of the rational number
func (r Rat) Numerator() int64 {
	return r.Num
}

// Denominator returns the denominator of the rational number
func (r Rat) Denominator() int64 {
	return r.Den
}
