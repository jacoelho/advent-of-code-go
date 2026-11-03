package matrix

import (
	"math"
	"slices"
	"testing"
)

func TestRREFUniqueSolutionAndClone(t *testing.T) {
	original, err := FromRows([][]Rat{
		{NewRat(0, 1), NewRat(2, 1), NewRat(8, 1)},
		{NewRat(3, 1), NewRat(1, 1), NewRat(7, 1)},
	})
	if err != nil {
		t.Fatal(err)
	}
	clone := original.Clone()
	result := RREF(clone)
	if result.Inconsistent || len(result.PivotCols) != 2 {
		t.Fatalf("expected unique solution, got %+v", result)
	}
	solution := ExtractParametricSolution(clone, result.PivotCols, 2)
	if len(solution.FreeCols) != 0 || solution.Base[0].Cmp(NewRat(1, 1)) != 0 || solution.Base[1].Cmp(NewRat(4, 1)) != 0 {
		t.Fatal("expected x=1, y=4")
	}
	for row, values := range [][]int64{{0, 2, 8}, {3, 1, 7}} {
		for col, want := range values {
			if original.Get(row, col).Cmp(NewRat(want, 1)) != 0 {
				t.Fatalf("RREF of clone changed original cell (%d,%d)", row, col)
			}
		}
	}
}

func TestRREFParametricAndInconsistent(t *testing.T) {
	for _, rightHandSide := range []int64{10, 11} {
		m, err := FromRows([][]Rat{
			{NewRat(1, 1), NewRat(2, 1), NewRat(5, 1)},
			{NewRat(2, 1), NewRat(4, 1), NewRat(rightHandSide, 1)},
		})
		if err != nil {
			t.Fatal(err)
		}
		result := RREF(m)
		if result.Inconsistent != (rightHandSide == 11) {
			t.Fatalf("right-hand side %d: got %+v", rightHandSide, result)
		}
		if result.Inconsistent {
			continue
		}
		solution := ExtractParametricSolution(m, result.PivotCols, 2)
		if len(result.PivotCols) != 1 || !slices.Equal(solution.FreeCols, []int{1}) {
			t.Fatalf("expected one free y component, got %+v", solution)
		}
		for _, parameter := range []int64{-3, 0, 4} {
			want := []int64{5 - 2*parameter, parameter}
			for i := range 2 {
				got := solution.Base[i].Add(solution.Coeffs[0][i].Mul(NewRat(parameter, 1)))
				if got.Cmp(NewRat(want[i], 1)) != 0 {
					t.Fatalf("parameter %d: component %d = %s/%s; want %d", parameter, i, got.Numerator(), got.Denominator(), want[i])
				}
			}
		}
	}
}

func TestScaleToIntegers(t *testing.T) {
	base, coefficients, denominator, err := ScaleToIntegers(
		[]Rat{NewRat(1, 2), NewRat(1, 3)},
		[][]Rat{{NewRat(-1, 6), NewRat(1, 1)}},
	)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(base, []int64{3, 2}) || len(coefficients) != 1 || !slices.Equal(coefficients[0], []int64{-1, 6}) || denominator != 6 {
		t.Fatalf("scaled solution = %v, %v, %d", base, coefficients, denominator)
	}
}

func TestScaleToIntegersOverflow(t *testing.T) {
	if _, _, _, err := ScaleToIntegers([]Rat{NewRat(1, math.MaxInt64), NewRat(1, 2)}, nil); err == nil {
		t.Fatal("expected common-denominator range error")
	}
	if _, _, _, err := ScaleToIntegers([]Rat{NewRat(1, 2)}, [][]Rat{{NewRat(math.MaxInt64, 1)}}); err == nil {
		t.Fatal("expected scaled-coefficient range error")
	}
}
