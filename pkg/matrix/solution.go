package matrix

import "github.com/jacoelho/advent-of-code-go/pkg/xmath"

// ParametricSolution represents the parametric solution to Ax = b.
// The solution is expressed as: x = Base + sum(FreeCols[i] * Coeffs[i])
type ParametricSolution struct {
	Base      []Rat   // base solution vector
	Coeffs    [][]Rat // coefficient matrix for free variables (each row is a free variable)
	FreeCols  []int   // indices of free variable columns
	PivotCols []int   // indices of pivot columns
}

// findPivotRowMapping maps each pivot column to its row index
func findPivotRowMapping(m *Matrix, varCount int) map[int]int {
	pivotForCol := make(map[int]int)
	for r := range m.data {
		row := m.data[r]
		for c := range varCount {
			if !row[c].IsZero() {
				pivotForCol[c] = r
				break
			}
		}
	}
	return pivotForCol
}

// findFreeColumns returns columns that are not pivot columns
func findFreeColumns(varCount int, pivotCols []int) []int {
	isPivot := make([]bool, varCount)
	for _, c := range pivotCols {
		isPivot[c] = true
	}
	freeCols := make([]int, 0, varCount-len(pivotCols))
	for c := range varCount {
		if !isPivot[c] {
			freeCols = append(freeCols, c)
		}
	}
	return freeCols
}

// ExtractParametricSolution extracts parametric solution from RREF matrix.
// The matrix m should be in RREF form (typically after calling RREF).
// varCount is the number of variables (excludes RHS column).
// pivotCols are the column indices that have pivots.
func ExtractParametricSolution(
	m *Matrix,
	pivotCols []int,
	varCount int,
) ParametricSolution {
	pivotForCol := findPivotRowMapping(m, varCount)
	freeCols := findFreeColumns(varCount, pivotCols)
	freeCount := len(freeCols)

	zero := NewRat(0, 1)
	base := make([]Rat, varCount)
	for i := range base {
		base[i] = zero
	}
	coeffs := New(freeCount, varCount)

	// set free variable coefficients to identity
	one := NewRat(1, 1)
	for idx, freeCol := range freeCols {
		coeffs.Set(idx, freeCol, one)
	}

	// extract base solution and pivot column coefficients
	rhsCol := m.Cols() - 1
	for _, pivotCol := range pivotCols {
		rowIdx := pivotForCol[pivotCol]
		base[pivotCol] = m.data[rowIdx][rhsCol]
		for fi, freeCol := range freeCols {
			coeffs.Set(fi, pivotCol, m.data[rowIdx][freeCol].Neg())
		}
	}

	return ParametricSolution{
		Base:      base,
		Coeffs:    coeffs.data,
		FreeCols:  freeCols,
		PivotCols: pivotCols,
	}
}

// ScaleToIntegers scales a parametric solution with rational coefficients to integer representation.
// Computes the LCM of all denominators and scales base and coefficients accordingly.
// Returns scaled base vector, scaled coefficient matrix, and the common denominator used.
func ScaleToIntegers(
	base []Rat,
	coeffs [][]Rat,
) ([]int64, [][]int64, int64) {
	buttonCount := len(base)
	freeCount := len(coeffs)

	commonDen := int64(1)
	for _, b := range base {
		commonDen = xmath.LCM(commonDen, b.Denominator())
	}
	for _, row := range coeffs {
		for _, v := range row {
			commonDen = xmath.LCM(commonDen, v.Denominator())
		}
	}

	scaledBase := make([]int64, buttonCount)
	for i, b := range base {
		scaledBase[i] = b.Numerator() * (commonDen / b.Denominator())
	}

	scaledCoeffs := make([][]int64, freeCount)
	for k := range freeCount {
		scaledCoeffs[k] = make([]int64, buttonCount)
		for i, v := range coeffs[k] {
			scaledCoeffs[k][i] = v.Numerator() * (commonDen / v.Denominator())
		}
	}

	return scaledBase, scaledCoeffs, commonDen
}
