package matrix

import (
	"fmt"
	"slices"
)

type Matrix struct {
	data [][]Rat
}

func New(rows, cols int) *Matrix {
	if rows < 0 || cols < 0 {
		panic("matrix dimensions must be non-negative")
	}
	mat := &Matrix{
		data: make([][]Rat, rows),
	}
	zero := NewRat(0, 1)
	for i := range rows {
		mat.data[i] = make([]Rat, cols)
		for j := range cols {
			mat.data[i][j] = zero
		}
	}
	return mat
}

func (m *Matrix) Set(row, col int, val Rat) {
	if row < 0 || row >= m.Rows() || col < 0 || col >= m.Cols() {
		panic("matrix index out of bounds")
	}
	m.data[row][col] = val
}

func (m *Matrix) Get(row, col int) Rat {
	if row < 0 || row >= m.Rows() || col < 0 || col >= m.Cols() {
		panic("matrix index out of bounds")
	}
	return m.data[row][col]
}

func (m *Matrix) Rows() int {
	return len(m.data)
}

func (m *Matrix) Cols() int {
	if len(m.data) == 0 {
		return 0
	}
	return len(m.data[0])
}

func (m *Matrix) Clone() *Matrix {
	result := &Matrix{data: make([][]Rat, len(m.data))}
	for i := range m.data {
		result.data[i] = slices.Clone(m.data[i])
	}
	return result
}

func (m *Matrix) Row(i int) []Rat {
	if i < 0 || i >= m.Rows() {
		panic("matrix index out of bounds")
	}
	return slices.Clone(m.data[i])
}

func FromRows(rows [][]Rat) (*Matrix, error) {
	if len(rows) == 0 {
		return &Matrix{data: [][]Rat{}}, nil
	}
	cols := len(rows[0])
	for i := 1; i < len(rows); i++ {
		if len(rows[i]) != cols {
			return nil, &ErrNonRectangular{Row: i, ExpectedCols: cols, ActualCols: len(rows[i])}
		}
	}
	mat := &Matrix{
		data: make([][]Rat, len(rows)),
	}
	for i, row := range rows {
		mat.data[i] = slices.Clone(row)
	}
	return mat, nil
}

func Identity(n int) *Matrix {
	if n < 0 {
		panic("matrix dimensions must be non-negative")
	}
	mat := New(n, n)
	one := NewRat(1, 1)
	for i := range n {
		mat.data[i][i] = one
	}
	return mat
}

type ErrNonRectangular struct {
	Row          int
	ExpectedCols int
	ActualCols   int
}

func (e *ErrNonRectangular) Error() string {
	return fmt.Sprintf("non-rectangular matrix: row %d has %d columns, expected %d", e.Row, e.ActualCols, e.ExpectedCols)
}
