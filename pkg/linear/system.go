// Package linear builds and solves exact rational linear systems.
//
// Declare unknowns with Variable and NewSystem, compose immutable expressions
// with Var.Scale and Expr.Add, and record equations with System.Add. Solve
// identifies free variables and expresses every declared variable in terms of
// them; callers never handle pivot or column indexes. A free variable is an
// original unknown whose value parameterises the solution family, not an
// unknown marked free before solving.
//
// Rat is the arithmetic type throughout. Scale and Equate convert signed
// integers exactly; ScaleRat and EquateRat accept computed or
// arbitrary-precision rationals.
//
// Elimination is dense and uses arbitrary-precision arithmetic, so its cost
// grows with variable count, equation count and the bit length of
// intermediate values.
package linear

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
)

// ErrInconsistent means the equations have no rational solution.
var ErrInconsistent = errors.New("inconsistent linear equations")

// System collects equations over a fixed set of declared variables.
// Its zero value is a system with no variables. Add mutates the builder;
// builders must not be copied after use or mutated concurrently.
// Callers bound equation counts, variable counts and input sizes before building.
type System[K Key] struct {
	variables []Var[K] // Never mutated; solutions share it.
	equations []Equation[K]
	err       error
}

// NewSystem snapshots the declarations. Duplicate keys are reported by Solve.
func NewSystem[K Key](variables ...Var[K]) *System[K] {
	for column, variable := range variables {
		if slices.Contains(variables[:column], variable) {
			return &System[K]{err: fmt.Errorf("duplicate variable %v", variable.key)}
		}
	}
	return &System[K]{variables: slices.Clone(variables)}
}

// Add records equation and returns s.
func (s *System[K]) Add(equation Equation[K]) *System[K] {
	s.equations = append(s.equations, equation)
	return s
}

func (s *System[K]) fillRow(row []Rat, equation Equation[K]) error {
	for _, term := range equation.terms {
		column := slices.Index(s.variables, term.variable)
		if column < 0 {
			return fmt.Errorf("undeclared variable %v", term.variable.key)
		}
		row[column] = term.coefficient
	}
	row[len(s.variables)] = equation.rhs
	return nil
}

// Solve leaves s unchanged and returns an independent solution.
// Free variables are original variables chosen by declaration-order pivoting.
// A system without equations leaves every declared variable free.
// Solve reports duplicate declarations first, then the first equation in
// insertion order that references an undeclared variable, then ErrInconsistent
// when no rational solution exists.
func (s *System[K]) Solve() (*Solution[K], error) {
	if s.err != nil {
		return nil, s.err
	}
	n := len(s.variables)
	width := n + 1
	data := make([]Rat, len(s.equations)*width)
	m := make([][]Rat, len(s.equations))
	for i, equation := range s.equations {
		m[i] = data[i*width : (i+1)*width : (i+1)*width]
		if err := s.fillRow(m[i], equation); err != nil {
			return nil, fmt.Errorf("equation %d: %w", i+1, err)
		}
	}
	pivotCols, consistent := rref(m)
	if !consistent {
		return nil, ErrInconsistent
	}
	pivotRows := make([]int, n)
	for column := range pivotRows {
		pivotRows[column] = -1
	}
	for row, column := range pivotCols {
		pivotRows[column] = row
	}
	var free []Var[K]
	var freeColumns []int
	for column, variable := range s.variables {
		if pivotRows[column] < 0 {
			free = append(free, variable)
			freeColumns = append(freeColumns, column)
		}
	}
	// Expression terms use key order; public free variables keep declaration order.
	slices.SortFunc(freeColumns, func(a, b int) int {
		return cmp.Compare(s.variables[a].key, s.variables[b].key)
	})
	expressions := make([]Expr[K], n)
	for column, variable := range s.variables {
		row := pivotRows[column]
		if row < 0 {
			expressions[column] = variable.Scale(1)
			continue
		}
		expression := Expr[K]{constant: m[row][n]}
		for _, freeColumn := range freeColumns {
			if coefficient := m[row][freeColumn]; !coefficient.IsZero() {
				expression.terms = append(expression.terms, term[K]{
					variable: s.variables[freeColumn], coefficient: coefficient.Neg(),
				})
			}
		}
		expressions[column] = expression
	}
	return &Solution[K]{variables: s.variables, free: free, expressions: expressions}, nil
}
