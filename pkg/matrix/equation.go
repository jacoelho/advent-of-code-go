package matrix

import "cmp"

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

func Variable[K Key](key K) Var[K] { return Var[K]{key: key} }

func (v Var[K]) Key() K { return v.key }

func (v Var[K]) Scale[N Signed](coefficient N) Expr[K] {
	return Expr[K]{coefficients: singleIntegerTerm(v, int64(coefficient))}
}

func (v Var[K]) ScaleRat(coefficient Rat) Expr[K] {
	return Expr[K]{coefficients: singleTerm(v, coefficient)}
}

type term[K Key] struct {
	variable    Var[K]
	coefficient Rat
}

// Expr is an immutable affine expression. Its zero value is zero.
// Only nonzero coefficients are retained; canceled variables are not references.
type Expr[K Key] struct {
	constant     Rat
	coefficients termSet[K]
}

func (e Expr[K]) Constant() Rat { return e.constant }

// Coefficient returns zero when variable is absent.
func (e Expr[K]) Coefficient(variable Var[K]) Rat {
	low, high := 0, e.coefficients.len()
	for low < high {
		middle := low + (high-low)/2
		candidate := e.coefficients.at(middle)
		switch cmp.Compare(candidate.variable.key, variable.key) {
		case -1:
			low = middle + 1
		case 1:
			high = middle
		case 0:
			return candidate.coefficient
		}
	}
	return Rat{}
}

func (e Expr[K]) Add(other Expr[K]) Expr[K] {
	if e.constant.IsZero() {
		e.constant = other.constant
	} else if !other.constant.IsZero() {
		e.constant = e.constant.Add(other.constant)
	}
	// Published terms are immutable; constant-only additions can share spill storage.
	leftCount, rightCount := e.coefficients.len(), other.coefficients.len()
	if leftCount == 0 {
		e.coefficients = other.coefficients
		return e
	}
	if rightCount == 0 {
		return e
	}
	if e.coefficients.spill == nil && other.coefficients.spill == nil &&
		leftCount == 1 && rightCount == 1 {
		// A value receiver owns its inline array; no spill storage is changed here.
		left, right := e.coefficients.inline[0].variable.key, other.coefficients.inline[0].variable.key
		switch {
		case left < right:
			e.coefficients.inline[1] = other.coefficients.inline[0]
		case left > right:
			e.coefficients.inline[1] = e.coefficients.inline[0]
			e.coefficients.inline[0] = other.coefficients.inline[0]
		default:
			variable := e.coefficients.inline[0].variable
			combined := integerRat(e.coefficients.inline[0].coefficient).Add(integerRat(other.coefficients.inline[0].coefficient))
			e.coefficients.inline[0] = integerTerm[K]{}
			if !combined.IsZero() {
				e.coefficients.append(term[K]{variable: variable, coefficient: combined})
			}
		}
		return e
	}
	e.coefficients = mergeTerms(&e.coefficients, &other.coefficients)
	return e
}

func (e Expr[K]) Equal[N Signed](rhs N) Equation[K] {
	return e.EqualRat(integerRat(int64(rhs)))
}

func (e Expr[K]) EqualRat(rhs Rat) Equation[K] {
	return Equation[K]{coefficients: e.coefficients, rhs: rhs.Sub(e.constant)}
}

// Equation is a complete equality. Its zero value represents 0 = 0.
type Equation[K Key] struct {
	coefficients termSet[K]
	rhs          Rat
}
