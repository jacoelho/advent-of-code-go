package matrix

import "github.com/jacoelho/advent-of-code-go/pkg/xmath"

// Rat represents a rational number using integers stored in reduced form (normalized)
type Rat struct {
	numerator, denominator int64
}

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
	return Rat{numerator: num, denominator: den}
}

func (r Rat) Add(other Rat) Rat {
	if r.denominator == other.denominator {
		return NewRat(r.numerator+other.numerator, r.denominator)
	}
	lcm := xmath.LCM(r.denominator, other.denominator)
	num := r.numerator*(lcm/r.denominator) + other.numerator*(lcm/other.denominator)
	return NewRat(num, lcm)
}

func (r Rat) Sub(other Rat) Rat {
	if r.denominator == other.denominator {
		return NewRat(r.numerator-other.numerator, r.denominator)
	}
	lcm := xmath.LCM(r.denominator, other.denominator)
	num := r.numerator*(lcm/r.denominator) - other.numerator*(lcm/other.denominator)
	return NewRat(num, lcm)
}

func (r Rat) Mul(other Rat) Rat {
	return NewRat(r.numerator*other.numerator, r.denominator*other.denominator)
}

func (r Rat) Div(other Rat) Rat {
	if other.numerator == 0 {
		panic("division by zero")
	}
	return NewRat(r.numerator*other.denominator, r.denominator*other.numerator)
}

func (r Rat) Neg() Rat {
	return Rat{numerator: -r.numerator, denominator: r.denominator}
}

func (r Rat) IsZero() bool {
	return r.numerator == 0
}
func (r Rat) Numerator() int64 {
	return r.numerator
}

func (r Rat) Denominator() int64 {
	return r.denominator
}
