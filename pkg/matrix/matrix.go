package matrix

// Matrix represents a matrix over a Field type
type Matrix[T Field[T]] struct {
	data [][]T
}

// New creates a new matrix with given dimensions.
// All elements are initialized to the Field zero value.
func New[T Field[T]](rows, cols int) *Matrix[T] {
	if rows < 0 || cols < 0 {
		panic("matrix dimensions must be non-negative")
	}
	mat := &Matrix[T]{
		data: make([][]T, rows),
	}
	var zero T
	zero = zero.Zero()
	for i := range rows {
		mat.data[i] = make([]T, cols)
		for j := range cols {
			mat.data[i][j] = zero
		}
	}
	return mat
}

// Set sets element at (row, col)
func (m *Matrix[T]) Set(row, col int, val T) {
	if row < 0 || row >= m.Rows() || col < 0 || col >= m.Cols() {
		panic("matrix index out of bounds")
	}
	m.data[row][col] = val
}

// Get gets element at (row, col)
func (m *Matrix[T]) Get(row, col int) T {
	if row < 0 || row >= m.Rows() || col < 0 || col >= m.Cols() {
		panic("matrix index out of bounds")
	}
	return m.data[row][col]
}

// Rows returns number of rows
func (m *Matrix[T]) Rows() int {
	return len(m.data)
}

// Cols returns number of columns
func (m *Matrix[T]) Cols() int {
	if len(m.data) == 0 {
		return 0
	}
	return len(m.data[0])
}

// Data returns the underlying matrix data (for internal use)
func (m *Matrix[T]) Data() [][]T {
	return m.data
}

// BuildFromBitmasks creates a coefficient matrix from bitmasks.
// Each mask defines a column: bit i set in masks[j] means row i, col j = 1.
// The last column contains rhs values.
func BuildFromBitmasks(masks []uint, rhs []int) *Matrix[Rat] {
	if len(masks) == 0 || len(rhs) == 0 {
		panic("masks and rhs must be non-empty")
	}

	rows := len(rhs)
	cols := len(masks) + 1 // +1 for RHS column

	mat := New[Rat](rows, cols)

	// Set RHS column
	for i := range rows {
		mat.Set(i, cols-1, NewRat(int64(rhs[i]), 1))
	}

	// Set coefficient columns from masks
	for j, mask := range masks {
		for i := range rows {
			if mask&(1<<i) != 0 {
				mat.Set(i, j, NewRat(1, 1))
			}
		}
	}

	return mat
}
