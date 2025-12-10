package matrix

import (
	"testing"
)

// testField is a simple integer field for testing
type testField int

func (t testField) Add(other testField) testField {
	return t + other
}

func (t testField) Sub(other testField) testField {
	return t - other
}

func (t testField) Mul(other testField) testField {
	return t * other
}

func (t testField) Div(other testField) testField {
	if other == 0 {
		panic("division by zero")
	}
	return t / other
}

func (t testField) Neg() testField {
	return -t
}

func (t testField) IsZero() bool {
	return t == 0
}

func (t testField) Zero() testField {
	return 0
}

func (t testField) One() testField {
	return 1
}

func TestMatrixNew(t *testing.T) {
	m := New[testField](3, 4)
	if m.Rows() != 3 {
		t.Errorf("expected 3 rows, got %d", m.Rows())
	}
	if m.Cols() != 4 {
		t.Errorf("expected 4 cols, got %d", m.Cols())
	}

	// Check all elements are zero
	for i := 0; i < 3; i++ {
		for j := 0; j < 4; j++ {
			if !m.Get(i, j).IsZero() {
				t.Errorf("element at (%d, %d) should be zero", i, j)
			}
		}
	}
}

func TestMatrixSetGet(t *testing.T) {
	m := New[testField](2, 2)
	m.Set(0, 0, 5)
	m.Set(0, 1, 10)
	m.Set(1, 0, 15)
	m.Set(1, 1, 20)

	if m.Get(0, 0) != 5 {
		t.Errorf("expected 5, got %d", m.Get(0, 0))
	}
	if m.Get(0, 1) != 10 {
		t.Errorf("expected 10, got %d", m.Get(0, 1))
	}
	if m.Get(1, 0) != 15 {
		t.Errorf("expected 15, got %d", m.Get(1, 0))
	}
	if m.Get(1, 1) != 20 {
		t.Errorf("expected 20, got %d", m.Get(1, 1))
	}
}

func TestRREF(t *testing.T) {
	// Test a simple 2x2 system: x + y = 3, 2x - y = 0
	// Augmented matrix: [1 1 | 3]
	//                   [2 -1 | 0]
	m := New[testField](2, 3)
	m.Set(0, 0, 1)
	m.Set(0, 1, 1)
	m.Set(0, 2, 3)
	m.Set(1, 0, 2)
	m.Set(1, 1, -1)
	m.Set(1, 2, 0)

	result := RREF(m)

	if result.Inconsistent {
		t.Error("system should be consistent")
	}

	if len(result.PivotCols) != 2 {
		t.Errorf("expected 2 pivot columns, got %d", len(result.PivotCols))
	}

	// After RREF, first row should be [1 0 | 1] and second [0 1 | 2]
	if m.Get(0, 0) != 1 || m.Get(0, 1) != 0 || m.Get(0, 2) != 1 {
		t.Errorf("first row should be [1 0 1], got [%d %d %d]", m.Get(0, 0), m.Get(0, 1), m.Get(0, 2))
	}
	if m.Get(1, 0) != 0 || m.Get(1, 1) != 1 || m.Get(1, 2) != 2 {
		t.Errorf("second row should be [0 1 2], got [%d %d %d]", m.Get(1, 0), m.Get(1, 1), m.Get(1, 2))
	}
}

func TestRREFInconsistent(t *testing.T) {
	// Test inconsistent system: x + y = 3, x + y = 5
	m := New[testField](2, 3)
	m.Set(0, 0, 1)
	m.Set(0, 1, 1)
	m.Set(0, 2, 3)
	m.Set(1, 0, 1)
	m.Set(1, 1, 1)
	m.Set(1, 2, 5)

	result := RREF(m)

	if !result.Inconsistent {
		t.Error("system should be inconsistent")
	}
}

func TestExtractParametricSolution(t *testing.T) {
	// Test system with free variable: x + y + z = 3, y = 1
	// After RREF: [1 0 1 | 2]
	//             [0 1 0 | 1]
	m := New[testField](2, 4)
	m.Set(0, 0, 1)
	m.Set(0, 1, 0)
	m.Set(0, 2, 1)
	m.Set(0, 3, 2)
	m.Set(1, 0, 0)
	m.Set(1, 1, 1)
	m.Set(1, 2, 0)
	m.Set(1, 3, 1)

	pivotCols := []int{0, 1}
	sol := ExtractParametricSolution(m, pivotCols, 3)

	if len(sol.FreeCols) != 1 {
		t.Errorf("expected 1 free column, got %d", len(sol.FreeCols))
	}
	if sol.FreeCols[0] != 2 {
		t.Errorf("expected free column 2, got %d", sol.FreeCols[0])
	}
	if len(sol.PivotCols) != 2 {
		t.Errorf("expected 2 pivot columns, got %d", len(sol.PivotCols))
	}
}

func TestMatrixEmpty(t *testing.T) {
	m := New[testField](0, 0)
	if m.Rows() != 0 {
		t.Errorf("expected 0 rows, got %d", m.Rows())
	}
	if m.Cols() != 0 {
		t.Errorf("expected 0 cols, got %d", m.Cols())
	}
}

func TestMatrix1x1(t *testing.T) {
	m := New[testField](1, 1)
	if m.Rows() != 1 {
		t.Errorf("expected 1 row, got %d", m.Rows())
	}
	if m.Cols() != 1 {
		t.Errorf("expected 1 col, got %d", m.Cols())
	}
	m.Set(0, 0, 42)
	if m.Get(0, 0) != 42 {
		t.Errorf("expected 42, got %d", m.Get(0, 0))
	}
}

func TestMatrixPanics(t *testing.T) {
	t.Run("negative dimensions", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative dimensions")
			}
		}()
		New[testField](-1, 5)
	})

	t.Run("out of bounds get", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for out of bounds access")
			}
		}()
		m := New[testField](2, 2)
		m.Get(5, 5)
	})

	t.Run("out of bounds set", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for out of bounds set")
			}
		}()
		m := New[testField](2, 2)
		m.Set(5, 5, 10)
	})
}

func TestRatArithmetic(t *testing.T) {
	// Test basic rational arithmetic
	a := NewRat(1, 2) // 1/2
	b := NewRat(1, 3) // 1/3

	// Addition: 1/2 + 1/3 = 5/6
	sum := a.Add(b)
	if sum.Num != 5 || sum.Den != 6 {
		t.Errorf("expected 5/6, got %d/%d", sum.Num, sum.Den)
	}

	// Subtraction: 1/2 - 1/3 = 1/6
	diff := a.Sub(b)
	if diff.Num != 1 || diff.Den != 6 {
		t.Errorf("expected 1/6, got %d/%d", diff.Num, diff.Den)
	}

	// Multiplication: 1/2 * 1/3 = 1/6
	prod := a.Mul(b)
	if prod.Num != 1 || prod.Den != 6 {
		t.Errorf("expected 1/6, got %d/%d", prod.Num, prod.Den)
	}

	// Division: 1/2 / 1/3 = 3/2
	quot := a.Div(b)
	if quot.Num != 3 || quot.Den != 2 {
		t.Errorf("expected 3/2, got %d/%d", quot.Num, quot.Den)
	}
}

func TestRatNormalization(t *testing.T) {
	// Test that rationals are normalized
	r := NewRat(6, 9) // Should be reduced to 2/3
	if r.Num != 2 || r.Den != 3 {
		t.Errorf("expected 2/3, got %d/%d", r.Num, r.Den)
	}

	// Test negative denominator handling
	r = NewRat(1, -2) // Should become -1/2
	if r.Num != -1 || r.Den != 2 {
		t.Errorf("expected -1/2, got %d/%d", r.Num, r.Den)
	}
}

func TestRatZeroDenominator(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for zero denominator")
		}
	}()
	NewRat(1, 0)
}

func TestRatMatrix(t *testing.T) {
	// Test using actual Rat type with matrix
	m := New[Rat](2, 2)
	m.Set(0, 0, NewRat(1, 2))
	m.Set(0, 1, NewRat(1, 3))
	m.Set(1, 0, NewRat(1, 4))
	m.Set(1, 1, NewRat(1, 5))

	if m.Get(0, 0).Num != 1 || m.Get(0, 0).Den != 2 {
		t.Errorf("expected 1/2, got %d/%d", m.Get(0, 0).Num, m.Get(0, 0).Den)
	}
}

func TestScaleToIntegers(t *testing.T) {
	// Test scaling rational solution to integers
	base := []Rat{NewRat(1, 2), NewRat(1, 3)}
	coeffs := [][]Rat{
		{NewRat(1, 4), NewRat(1, 6)},
	}

	scaledBase, scaledCoeffs, commonDen := ScaleToIntegers(base, coeffs)

	// Common denominator should be LCM(2, 3, 4, 6) = 12
	if commonDen != 12 {
		t.Errorf("expected common denominator 12, got %d", commonDen)
	}

	// Scaled base: [6, 4]
	if scaledBase[0] != 6 || scaledBase[1] != 4 {
		t.Errorf("expected [6, 4], got [%d, %d]", scaledBase[0], scaledBase[1])
	}

	// Scaled coeffs: [[3, 2]]
	if scaledCoeffs[0][0] != 3 || scaledCoeffs[0][1] != 2 {
		t.Errorf("expected [[3, 2]], got [[%d, %d]]", scaledCoeffs[0][0], scaledCoeffs[0][1])
	}
}

func TestBuildFromBitmasks(t *testing.T) {
	masks := []uint{0b101, 0b110} // First affects rows 0,2; second affects rows 1,2
	rhs := []int{5, 7, 9}

	mat := BuildFromBitmasks(masks, rhs)

	// Verify matrix dimensions
	if mat.Rows() != 3 {
		t.Errorf("expected 3 rows, got %d", mat.Rows())
	}
	if mat.Cols() != 3 {
		t.Errorf("expected 3 cols, got %d", mat.Cols())
	}

	// Verify coefficients
	// Row 0: [1, 0, 5]
	if mat.Get(0, 0) != NewRat(1, 1) {
		t.Errorf("expected 1/1, got %d/%d", mat.Get(0, 0).Num, mat.Get(0, 0).Den)
	}
	if mat.Get(0, 1) != NewRat(0, 1) {
		t.Errorf("expected 0/1, got %d/%d", mat.Get(0, 1).Num, mat.Get(0, 1).Den)
	}
	if mat.Get(0, 2) != NewRat(5, 1) {
		t.Errorf("expected 5/1, got %d/%d", mat.Get(0, 2).Num, mat.Get(0, 2).Den)
	}

	// Row 1: [0, 1, 7]
	if mat.Get(1, 0) != NewRat(0, 1) {
		t.Errorf("expected 0/1, got %d/%d", mat.Get(1, 0).Num, mat.Get(1, 0).Den)
	}
	if mat.Get(1, 1) != NewRat(1, 1) {
		t.Errorf("expected 1/1, got %d/%d", mat.Get(1, 1).Num, mat.Get(1, 1).Den)
	}
	if mat.Get(1, 2) != NewRat(7, 1) {
		t.Errorf("expected 7/1, got %d/%d", mat.Get(1, 2).Num, mat.Get(1, 2).Den)
	}

	// Row 2: [1, 1, 9]
	if mat.Get(2, 0) != NewRat(1, 1) {
		t.Errorf("expected 1/1, got %d/%d", mat.Get(2, 0).Num, mat.Get(2, 0).Den)
	}
	if mat.Get(2, 1) != NewRat(1, 1) {
		t.Errorf("expected 1/1, got %d/%d", mat.Get(2, 1).Num, mat.Get(2, 1).Den)
	}
	if mat.Get(2, 2) != NewRat(9, 1) {
		t.Errorf("expected 9/1, got %d/%d", mat.Get(2, 2).Num, mat.Get(2, 2).Den)
	}
}

func TestBuildFromBitmasksPanic(t *testing.T) {
	t.Run("empty masks", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for empty masks")
			}
		}()
		BuildFromBitmasks([]uint{}, []int{1, 2})
	})

	t.Run("empty rhs", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for empty rhs")
			}
		}()
		BuildFromBitmasks([]uint{1}, []int{})
	})
}
