package matrix

type rrefResult struct {
	PivotCols    []int
	Inconsistent bool
}

func findPivotRow(m [][]Rat, startRow, col int) int {
	for i := startRow; i < len(m); i++ {
		if !m[i][col].IsZero() {
			return i
		}
	}
	return -1
}

func normalizePivotRow(row []Rat, pivotCol int) {
	pivotVal := row[pivotCol]
	if pivotVal.IsZero() {
		panic("cannot normalize row with zero pivot")
	}
	for j := pivotCol; j < len(row); j++ {
		row[j] = row[j].Div(pivotVal)
	}
}

func eliminateColumn(m [][]Rat, pivotRow, col int) {
	pivotRowData := m[pivotRow]
	cols := len(m[0])
	for i := range m {
		if i == pivotRow {
			continue
		}
		factor := m[i][col]
		if factor.IsZero() {
			continue
		}
		row := m[i]
		for j := col; j < cols; j++ {
			row[j] = row[j].Sub(factor.Mul(pivotRowData[j]))
		}
	}
}

// rref reduces augmented matrix [A|b] in place; the last column is the right-hand side.
// PivotCols lists variable columns in pivot-row order. Inconsistent reports a row 0=b with b nonzero.
func rref(m [][]Rat) rrefResult {
	rows := len(m)
	if rows == 0 {
		return rrefResult{}
	}
	cols := len(m[0]) - 1

	pivotCols := make([]int, 0, min(rows, cols))
	r := 0

	for c := 0; c < cols && r < rows; c++ {
		pivot := findPivotRow(m, r, c)
		if pivot == -1 {
			continue
		}

		m[r], m[pivot] = m[pivot], m[r]
		normalizePivotRow(m[r], c)
		eliminateColumn(m, r, c)

		pivotCols = append(pivotCols, c)
		r++
	}

	// Elimination leaves non-pivot rows with zero coefficients.
	for _, row := range m[r:] {
		if !row[cols].IsZero() {
			return rrefResult{PivotCols: pivotCols, Inconsistent: true}
		}
	}

	return rrefResult{PivotCols: pivotCols, Inconsistent: false}
}
