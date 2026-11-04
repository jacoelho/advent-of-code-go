package matrix

import (
	"math"
	"slices"
	"strings"
	"testing"
)

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
	half, _ := solveDiagonalRats(t, 2, 1, 3, 1)
	for _, tt := range []struct {
		name         string
		base         []Rat
		coefficients [][]Rat
		context      string
	}{
		{"common denominator", []Rat{NewRat(1, math.MaxInt64), half}, nil, "common denominator exceeds int64"},
		{"scaled coefficient", []Rat{half}, [][]Rat{{NewRat(math.MaxInt64, 1)}}, "scale coefficient row 0: scaled value 0 exceeds int64"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			base, coefficients, denominator, err := ScaleToIntegers(tt.base, tt.coefficients)
			if err == nil || !strings.Contains(err.Error(), tt.context) || base != nil || coefficients != nil || denominator != 0 {
				t.Fatalf("got %v, %v, %d, %v; want error containing %q without partial values", base, coefficients, denominator, err, tt.context)
			}
		})
	}
}

func TestScaleSolvedFractions(t *testing.T) {
	half, third := solveDiagonalRats(t, 2, 1, 3, 1)
	base, coefficients, denominator, err := ScaleToIntegers([]Rat{half}, [][]Rat{{third}})
	if err != nil || !slices.Equal(base, []int64{3}) || len(coefficients) != 1 || !slices.Equal(coefficients[0], []int64{2}) || denominator != 6 {
		t.Fatalf("scaled fractions = %v, %v, %d, %v", base, coefficients, denominator, err)
	}
	largeHalf, _ := solveDiagonalRats(t, 2, 1, 3_000_000_000_000_000_001, 0)
	base, _, denominator, err = ScaleToIntegers([]Rat{largeHalf, NewRat(1, 4)}, nil)
	if err != nil || !slices.Equal(base, []int64{2, 1}) || denominator != 4 {
		t.Fatalf("scaled reducible fraction = %v, denominator %d, error %v; want [2 1], 4, nil", base, denominator, err)
	}
}

func TestScaleToIntegersFractionBoundaries(t *testing.T) {
	for _, tt := range []struct {
		name        string
		values      []Rat
		want        []int64
		denominator int64
	}{
		{"LCM avoids product overflow", []Rat{NewRat(1, 4_000_000_000), NewRat(1, 6_000_000_000)}, []int64{3, 2}, 12_000_000_000},
		{"minimum numerator", []Rat{NewRat(math.MinInt64, 3), NewRat(1, 3)}, []int64{math.MinInt64, 1}, 3},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, _, denominator, err := ScaleToIntegers(tt.values, nil)
			if err != nil || !slices.Equal(got, tt.want) || denominator != tt.denominator {
				t.Fatalf("scaled values = %v, denominator %d, error %v; want %v, %d", got, denominator, err, tt.want, tt.denominator)
			}
		})
	}
	for _, tt := range []struct {
		name    string
		values  []Rat
		context string
	}{
		{"scaled numerator underflow", []Rat{NewRat(math.MinInt64, 3), NewRat(1, 6)}, "scale base: scaled value 0 exceeds int64"},
		{"denominator exceeds positive range", []Rat{NewRat(1, math.MinInt64)}, "common denominator exceeds int64"},
		{"numerator exceeds positive range", []Rat{NewRat(math.MinInt64, -1)}, "scale base: scaled value 0 exceeds int64"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			values, coefficients, denominator, err := ScaleToIntegers(tt.values, nil)
			if err == nil || !strings.Contains(err.Error(), tt.context) || values != nil || coefficients != nil || denominator != 0 {
				t.Fatalf("got %v, %v, %d, %v; want error containing %q without partial values", values, coefficients, denominator, err, tt.context)
			}
		})
	}
}
