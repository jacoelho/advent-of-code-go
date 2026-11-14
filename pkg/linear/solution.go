package linear

import (
	"fmt"
	"slices"
)

// Solution is an immutable rational solution family. Its zero value is the
// solution of an empty system. Integrality, bounds and objectives are caller policy.
type Solution[K Key] struct {
	variables   []Var[K]
	free        []Var[K]
	expressions []Expr[K] // Indexed by declaration order.
}

// Free returns an independent slice in declared variable order.
// An empty slice means the solution is unique.
func (s *Solution[K]) Free() []Var[K] { return slices.Clone(s.free) }

// Expression returns variable expressed in terms of Free. It returns an error
// if variable was not declared.
func (s *Solution[K]) Expression(variable Var[K]) (Expr[K], error) {
	column := slices.Index(s.variables, variable)
	if column < 0 {
		return Expr[K]{}, fmt.Errorf("undeclared variable %v", variable.key)
	}
	return s.expressions[column], nil
}

// Bind requires exactly one assignment per free variable and returns all
// declared values in an independent map. Failure returns no partial values.
func (s *Solution[K]) Bind(assignments map[Var[K]]Rat) (map[Var[K]]Rat, error) {
	for _, variable := range s.free {
		if _, assigned := assignments[variable]; !assigned {
			return nil, fmt.Errorf("missing free variable %v", variable.key)
		}
	}
	if len(assignments) != len(s.free) {
		return nil, fmt.Errorf("assignments must contain only the %d free variables", len(s.free))
	}
	values := make(map[Var[K]]Rat, len(s.variables))
	for column, variable := range s.variables {
		expression := s.expressions[column]
		value := expression.constant
		for _, term := range expression.terms {
			value = value.Add(term.coefficient.Mul(assignments[term.variable]))
		}
		values[variable] = value
	}
	return values, nil
}
