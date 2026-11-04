package matrix

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
)

// ErrInconsistent means the equations have no rational solution.
var ErrInconsistent = errors.New("inconsistent linear equations")

// declarations owns its inline values or shares a constructor-owned immutable slice.
type declarations[K Key] struct {
	inline [2]Var[K]
	count  uint8
	spill  []Var[K]
}

func (d *declarations[K]) len() int {
	if d.spill != nil {
		return len(d.spill)
	}
	return int(d.count)
}

func (d *declarations[K]) at(index int) Var[K] {
	if d.spill != nil {
		return d.spill[index]
	}
	return d.inline[index]
}

func (d *declarations[K]) index(variable Var[K]) int {
	if d.spill != nil {
		return slices.Index(d.spill, variable)
	}
	for i := range int(d.count) {
		if d.inline[i] == variable {
			return i
		}
	}
	return -1
}

func declareVariables[K Key](variables []Var[K]) (declarations[K], error) {
	for column, variable := range variables {
		if slices.Contains(variables[:column], variable) {
			return declarations[K]{}, fmt.Errorf("duplicate variable %v", variable.key)
		}
	}
	if len(variables) > 2 {
		return declarations[K]{spill: slices.Clone(variables)}, nil
	}
	declared := declarations[K]{count: uint8(len(variables))}
	copy(declared.inline[:], variables)
	return declared, nil
}

// System builds exact rational equations in declared variable order.
// Its zero value is a system with no variables. Add mutates the builder;
// builders must not be copied after use or mutated concurrently.
// Callers bound equation counts, variable counts and input sizes before building.
type System[K Key] struct {
	variables declarations[K]
	inline    [2]Equation[K]
	equations []Equation[K] // Full owned sequence after the third equation.
	count     uint8
	err       error
}

// NewSystem snapshots the declarations. Duplicate keys are reported by Solve.
func NewSystem[K Key](variables ...Var[K]) *System[K] {
	declared, err := declareVariables(variables)
	return &System[K]{variables: declared, err: err}
}

func (s *System[K]) equationCount() int {
	if s.equations != nil {
		return len(s.equations)
	}
	return int(s.count)
}

func (s *System[K]) equation(index int) *Equation[K] {
	if s.equations != nil {
		return &s.equations[index]
	}
	return &s.inline[index]
}

// Add snapshots an equation. Solve validates variable references in insertion order.
// Later additions are retained even if an earlier equation is invalid.
func (s *System[K]) Add(equation Equation[K]) *System[K] {
	if s.err != nil {
		return s
	}
	if s.equations != nil {
		s.equations = append(s.equations, equation)
	} else if s.count < 2 {
		s.inline[s.count] = equation
		s.count++
	} else {
		s.equations = []Equation[K]{s.inline[0], s.inline[1], equation}
		clear(s.inline[:])
		s.count = 0
	}
	return s
}

func (s *System[K]) fillRow(row []Rat, equation *Equation[K]) error {
	for i := range equation.coefficients.len() {
		term := equation.coefficients.at(i)
		column := s.variables.index(term.variable)
		if column < 0 {
			return fmt.Errorf("undeclared variable %v", term.variable.key)
		}
		row[column] = term.coefficient
	}
	row[s.variables.len()] = equation.rhs
	return nil
}

// Solve preserves the builder and returns an independent solution snapshot.
// Free variables are original variables chosen by declaration-order pivoting.
// A system without equations leaves every declared variable free.
func (s *System[K]) Solve() (solution *Solution[K], err error) {
	var result Solution[K]
	err = s.solveInto(&result)
	if err == nil {
		solution = &result
	}
	return
}

// solveInto requires a fresh, zeroed result.
func (s *System[K]) solveInto(result *Solution[K]) error {
	if s.err != nil {
		return s.err
	}
	if s.variables.len() == 2 && s.equationCount() == 2 {
		var rows [2][3]Rat
		for i := range 2 {
			if err := s.fillRow(rows[i][:], s.equation(i)); err != nil {
				return fmt.Errorf("equation %d: %w", i+1, err)
			}
		}
		if x, y, ok := solveTwoIntegers(&rows); ok {
			result.variables = s.variables
			result.unique[0] = x
			result.unique[1] = y
			return nil
		}
	}
	width := s.variables.len() + 1
	data := make([]Rat, s.equationCount()*width)
	m := make([][]Rat, s.equationCount())
	for i := range s.equationCount() {
		m[i] = data[i*width : (i+1)*width : (i+1)*width]
		if err := s.fillRow(m[i], s.equation(i)); err != nil {
			return fmt.Errorf("equation %d: %w", i+1, err)
		}
	}
	reduced := rref(m)
	if reduced.Inconsistent {
		return ErrInconsistent
	}
	pivotRows := make([]int, s.variables.len())
	for column := range pivotRows {
		pivotRows[column] = -1
	}
	for row, column := range reduced.PivotCols {
		pivotRows[column] = row
	}
	result.variables = s.variables
	var freeColumns []int
	for column := range s.variables.len() {
		variable := s.variables.at(column)
		if pivotRows[column] < 0 {
			result.free = append(result.free, variable)
			freeColumns = append(freeColumns, column)
		}
	}
	if s.variables.len() == 2 && len(result.free) == 0 {
		result.unique[0] = m[pivotRows[0]][2]
		result.unique[1] = m[pivotRows[1]][2]
		return nil
	}
	result.expressions = make([]Expr[K], s.variables.len())
	// Expression terms use key order; public free variables keep declaration order.
	slices.SortFunc(freeColumns, func(a, b int) int {
		return cmp.Compare(s.variables.at(a).key, s.variables.at(b).key)
	})
	for column := range s.variables.len() {
		variable := s.variables.at(column)
		row := pivotRows[column]
		if row < 0 {
			result.expressions[column] = variable.Scale(1)
			continue
		}
		expression := Expr[K]{constant: m[row][s.variables.len()]}
		for _, freeColumn := range freeColumns {
			coefficient := m[row][freeColumn]
			if coefficient.IsZero() {
				continue
			}
			expression.coefficients.append(term[K]{
				variable: s.variables.at(freeColumn), coefficient: coefficient.Neg(),
			})
		}
		result.expressions[column] = expression
	}
	return nil
}

// Solution is an immutable rational solution family. Its zero value is the
// solution of an empty system. Integrality, bounds and objectives are caller policy.
type Solution[K Key] struct {
	variables   declarations[K]
	unique      [2]Rat // Canonical for exactly two variables with no free variables.
	free        []Var[K]
	expressions []Expr[K]
}

// Free returns an independent slice in declared variable order.
// An empty slice means the solution is unique.
func (s *Solution[K]) Free() []Var[K] { return slices.Clone(s.free) }

// Expression expresses variable in terms of Free. Undeclared keys are errors.
func (s *Solution[K]) Expression(variable Var[K]) (expression Expr[K], err error) {
	var coefficients *termSet[K]
	expression.constant, coefficients, err = s.lookupExpression(variable)
	if coefficients != nil {
		expression.coefficients = *coefficients
	}
	return
}

func (s *Solution[K]) lookupExpression(variable Var[K]) (Rat, *termSet[K], error) {
	column := s.variables.index(variable)
	if column < 0 {
		return Rat{}, nil, fmt.Errorf("undeclared variable %v", variable.key)
	}
	constant, coefficients := s.expressionParts(column)
	return constant, coefficients, nil
}

// expressionParts borrows immutable coefficients; nil denotes no terms.
// Public expressions copy the term set instead of retaining this pointer.
func (s *Solution[K]) expressionParts(column int) (Rat, *termSet[K]) {
	if s.variables.len() == 2 && len(s.free) == 0 {
		return s.unique[column], nil
	}
	expression := &s.expressions[column]
	return expression.constant, &expression.coefficients
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
	values := make(map[Var[K]]Rat, s.variables.len())
	for column := range s.variables.len() {
		variable := s.variables.at(column)
		value, coefficients := s.expressionParts(column)
		if coefficients != nil {
			for i := range coefficients.len() {
				term := coefficients.at(i)
				value = value.Add(term.coefficient.Mul(assignments[term.variable]))
			}
		}
		values[variable] = value
	}
	return values, nil
}
