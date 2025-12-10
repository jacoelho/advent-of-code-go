package matrix

// RREFResult contains the result of RREF operation
type RREFResult struct {
	PivotCols    []int // columns that have pivots
	Inconsistent bool  // true if system is inconsistent
}

// findPivotRow finds the first non-zero entry in a column starting from startRow
func findPivotRow[T Field[T]](m *Matrix[T], startRow, col int) int {
	for i := startRow; i < m.Rows(); i++ {
		if !m.data[i][col].IsZero() {
			return i
		}
	}
	return -1
}

// normalizePivotRow divides the pivot row by its pivot value to make the pivot 1
func normalizePivotRow[T Field[T]](row []T, pivotCol int) {
	pivotVal := row[pivotCol]
	if pivotVal.IsZero() {
		panic("cannot normalize row with zero pivot")
	}
	for j := pivotCol; j < len(row); j++ {
		row[j] = row[j].Div(pivotVal)
	}
}

// eliminateColumn eliminates the pivot column from all other rows
func eliminateColumn[T Field[T]](m *Matrix[T], pivotRow, col int) {
	for i := range m.data {
		if i == pivotRow {
			continue
		}
		factor := m.data[i][col]
		if factor.IsZero() {
			continue
		}
		for j := col; j < m.Cols(); j++ {
			m.data[i][j] = m.data[i][j].Sub(factor.Mul(m.data[pivotRow][j]))
		}
	}
}

// RREF performs Reduced Row Echelon Form on augmented matrix [A|b].
// Last column is treated as RHS (right-hand side).
// The matrix is modified in place.
func RREF[T Field[T]](m *Matrix[T]) RREFResult {
	rows := m.Rows()
	if rows == 0 {
		return RREFResult{}
	}
	cols := m.Cols() - 1 // last column is RHS

	pivotCols := make([]int, 0, min(rows, cols))
	r := 0

	for c := 0; c < cols && r < rows; c++ {
		pivot := findPivotRow(m, r, c)
		if pivot == -1 {
			continue
		}

		m.data[r], m.data[pivot] = m.data[pivot], m.data[r]
		normalizePivotRow(m.data[r], c)
		eliminateColumn(m, r, c)

		pivotCols = append(pivotCols, c)
		r++
	}

	// detect inconsistency (0 ... 0 | b) with b != 0
	for i := range rows {
		allZero := true
		for c := range cols {
			if !m.data[i][c].IsZero() {
				allZero = false
				break
			}
		}
		if allZero && !m.data[i][cols].IsZero() {
			return RREFResult{PivotCols: pivotCols, Inconsistent: true}
		}
	}

	return RREFResult{PivotCols: pivotCols, Inconsistent: false}
}
