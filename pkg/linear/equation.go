package linear

import (
	"cmp"
	"slices"
)

// Key identifies a variable by value. Strings and integers provide reflexive
// equality without the NaN and dynamic-comparability cases of arbitrary keys.
type Key interface {
	~string | ~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Signed permits exact conversion of integer coefficients to Rat.
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Var names an unknown. Equal keys identify the same variable across systems.
type Var[K Key] struct {
	key K
}

// Variable returns the unknown identified by key.
func Variable[K Key](key K) Var[K] { return Var[K]{key: key} }

// Key returns the key that identifies v.
func (v Var[K]) Key() K { return v.key }

// Scale returns the expression coefficient·v.
func (v Var[K]) Scale[N Signed](coefficient N) Expr[K] {
	return v.ScaleRat(integerRat(int64(coefficient)))
}

// ScaleRat returns the expression coefficient·v.
func (v Var[K]) ScaleRat(coefficient Rat) Expr[K] {
	if coefficient.IsZero() {
		return Expr[K]{}
	}
	return Expr[K]{terms: []term[K]{{variable: v, coefficient: coefficient}}}
}

// Expr is an immutable affine expression. Its zero value is zero.
// Only nonzero coefficients are retained, so a variable whose coefficients
// cancel is not checked against a system's declarations.
type Expr[K Key] struct {
	constant Rat
	terms    []term[K] // Sorted by key; never mutated after construction.
}

// Constant returns the term of e that multiplies no variable.
func (e Expr[K]) Constant() Rat { return e.constant }

// Coefficient returns the coefficient of variable in e, or zero when absent.
func (e Expr[K]) Coefficient(variable Var[K]) Rat {
	index, found := slices.BinarySearchFunc(e.terms, variable.key, func(t term[K], key K) int {
		return cmp.Compare(t.variable.key, key)
	})
	if !found {
		return Rat{}
	}
	return e.terms[index].coefficient
}

// Add returns e + other.
func (e Expr[K]) Add(other Expr[K]) Expr[K] {
	return Expr[K]{constant: e.constant.Add(other.constant), terms: mergeTerms(e.terms, other.terms)}
}

// Equate returns the equation e = rhs.
func (e Expr[K]) Equate[N Signed](rhs N) Equation[K] {
	return e.EquateRat(integerRat(int64(rhs)))
}

// EquateRat returns the equation e = rhs.
func (e Expr[K]) EquateRat(rhs Rat) Equation[K] {
	return Equation[K]{terms: e.terms, rhs: rhs.Sub(e.constant)}
}

// Equation is a complete equality. Its zero value represents 0 = 0.
type Equation[K Key] struct {
	terms []term[K]
	rhs   Rat
}
