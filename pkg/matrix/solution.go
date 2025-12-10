package matrix

import "github.com/jacoelho/advent-of-code-go/pkg/xmath"

// Rational interface for types that can be scaled to integers
type Rational[T any] interface {
	Field[T]
	Numerator() int64
	Denominator() int64
}

// ParametricSolution represents the parametric solution to Ax = b
type ParametricSolution[T Field[T]] struct {
	Base      []T   // base solution vector
	Coeffs    [][]T // coefficient matrix for free variables
	FreeCols  []int // indices of free variable columns
	PivotCols []int // indices of pivot columns
}

// findPivotRowMapping maps each pivot column to its row index
func findPivotRowMapping[T Field[T]](m *Matrix[T], varCount int) map[int]int {
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
// varCount is the number of variables (excludes RHS column).
func ExtractParametricSolution[T Field[T]](
	m *Matrix[T],
	pivotCols []int,
	varCount int,
) ParametricSolution[T] {
	pivotForCol := findPivotRowMapping(m, varCount)
	freeCols := findFreeColumns(varCount, pivotCols)
	freeCount := len(freeCols)

	var zero T
	zero = zero.Zero()
	base := make([]T, varCount)
	for i := range base {
		base[i] = zero
	}
	coeffs := New[T](freeCount, varCount)

	// Set free variable coefficients to identity
	var one T
	one = one.One()
	for idx, freeCol := range freeCols {
		coeffs.Set(idx, freeCol, one)
	}

	// Extract base solution and pivot column coefficients
	rhsCol := m.Cols() - 1
	for _, pivotCol := range pivotCols {
		rowIdx := pivotForCol[pivotCol]
		base[pivotCol] = m.data[rowIdx][rhsCol]
		for fi, freeCol := range freeCols {
			coeffs.Set(fi, pivotCol, m.data[rowIdx][freeCol].Neg())
		}
	}

	return ParametricSolution[T]{
		Base:      base,
		Coeffs:    coeffs.Data(),
		FreeCols:  freeCols,
		PivotCols: pivotCols,
	}
}

// ScaleToIntegers scales a parametric solution with rational coefficients to integer representation.
// Returns scaled base vector, scaled coefficient matrix, and common denominator.
func ScaleToIntegers[T Rational[T]](
	base []T,
	coeffs [][]T,
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
