package linear

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
	for j := pivotCol; j < len(row); j++ {
		row[j] = row[j].Div(pivotVal)
	}
}

func eliminateColumn(m [][]Rat, pivotRow, col int) {
	pivotRowData := m[pivotRow]
	for i, row := range m {
		if i == pivotRow {
			continue
		}
		factor := row[col]
		if factor.IsZero() {
			continue
		}
		for j := col; j < len(row); j++ {
			row[j] = row[j].Sub(factor.Mul(pivotRowData[j]))
		}
	}
}

// rref reduces augmented matrix [A|b] in place; the last column is the right-hand side.
// pivotCols lists variable columns in pivot-row order. consistent is false when
// a row reduces to 0 = b with b nonzero.
func rref(m [][]Rat) (pivotCols []int, consistent bool) {
	rows := len(m)
	if rows == 0 {
		return nil, true
	}
	cols := len(m[0]) - 1
	pivotCols = make([]int, 0, min(rows, cols))
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
			return pivotCols, false
		}
	}
	return pivotCols, true
}
