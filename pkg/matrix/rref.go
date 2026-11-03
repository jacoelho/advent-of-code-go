package matrix

type RREFResult struct {
	PivotCols    []int
	Inconsistent bool
}

func findPivotRow(m *Matrix, startRow, col int) int {
	for i := startRow; i < m.Rows(); i++ {
		if !m.data[i][col].IsZero() {
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

func eliminateColumn(m *Matrix, pivotRow, col int) {
	pivotRowData := m.data[pivotRow]
	cols := m.Cols()
	for i := range m.data {
		if i == pivotRow {
			continue
		}
		factor := m.data[i][col]
		if factor.IsZero() {
			continue
		}
		row := m.data[i]
		for j := col; j < cols; j++ {
			row[j] = row[j].Sub(factor.Mul(pivotRowData[j]))
		}
	}
}

// RREF reduces augmented matrix [A|b] in place; the last column is the right-hand side.
// PivotCols lists variable columns in pivot-row order. Inconsistent reports a row 0=b with b nonzero.
func RREF(m *Matrix) RREFResult {
	rows := m.Rows()
	if rows == 0 {
		return RREFResult{}
	}
	cols := m.Cols() - 1

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

	rhsCol := cols
	for i := range rows {
		row := m.data[i]
		allZero := true
		for c := range cols {
			if !row[c].IsZero() {
				allZero = false
				break
			}
		}
		if allZero && !row[rhsCol].IsZero() {
			return RREFResult{PivotCols: pivotCols, Inconsistent: true}
		}
	}

	return RREFResult{PivotCols: pivotCols, Inconsistent: false}
}
