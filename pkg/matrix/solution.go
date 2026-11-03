package matrix

import (
	"fmt"
	"math/big"
)

// ParametricSolution expresses x = Base + sum(parameters[i] * Coeffs[i]).
// FreeCols[i] identifies the component of x equal to parameters[i].
type ParametricSolution struct {
	Base      []Rat
	Coeffs    [][]Rat
	FreeCols  []int
	PivotCols []int
}

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

// ExtractParametricSolution returns the solutions of a consistent augmented
// matrix m already reduced by [RREF]. pivotCols must contain its pivot columns;
// varCount excludes the right-hand-side column.
func ExtractParametricSolution(
	m *Matrix,
	pivotCols []int,
	varCount int,
) ParametricSolution {
	freeCols := findFreeColumns(varCount, pivotCols)
	freeCount := len(freeCols)

	zero := NewRat(0, 1)
	base := make([]Rat, varCount)
	for i := range base {
		base[i] = zero
	}
	coeffs := New(freeCount, varCount)

	one := NewRat(1, 1)
	for idx, freeCol := range freeCols {
		coeffs.Set(idx, freeCol, one)
	}

	rhsCol := m.Cols() - 1
	for rowIdx, pivotCol := range pivotCols {
		base[pivotCol] = m.data[rowIdx][rhsCol]
		for fi, freeCol := range freeCols {
			coeffs.Set(fi, pivotCol, m.data[rowIdx][freeCol].Neg())
		}
	}

	return ParametricSolution{
		Base:      base,
		Coeffs:    coeffs.data,
		FreeCols:  freeCols,
		PivotCols: pivotCols,
	}
}

// ScaleToIntegers multiplies base and coeffs by one positive common
// denominator. It returns the scaled base, scaled coefficients and denominator,
// or an error if any scaled value or the denominator cannot fit in int64.
func ScaleToIntegers(base []Rat, coeffs [][]Rat) ([]int64, [][]int64, int64, error) {
	commonDen := big.NewInt(1)
	includeDenominator := func(value Rat) {
		denominator := value.rational().Denom()
		gcd := new(big.Int).GCD(nil, nil, commonDen, denominator)
		commonDen.Quo(commonDen, gcd)
		commonDen.Mul(commonDen, denominator)
	}
	for _, value := range base {
		includeDenominator(value)
	}
	for _, row := range coeffs {
		for _, value := range row {
			includeDenominator(value)
		}
	}
	if !commonDen.IsInt64() {
		return nil, nil, 0, fmt.Errorf("common denominator exceeds int64")
	}

	scale := func(values []Rat) ([]int64, error) {
		result := make([]int64, len(values))
		for i, value := range values {
			factor := new(big.Int).Quo(commonDen, value.rational().Denom())
			numerator := new(big.Int).Mul(value.rational().Num(), factor)
			if !numerator.IsInt64() {
				return nil, fmt.Errorf("scaled value %d exceeds int64", i)
			}
			result[i] = numerator.Int64()
		}
		return result, nil
	}
	scaledBase, err := scale(base)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("scale base: %w", err)
	}
	scaledCoeffs := make([][]int64, len(coeffs))
	for i, row := range coeffs {
		scaledCoeffs[i], err = scale(row)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("scale coefficient row %d: %w", i, err)
		}
	}
	return scaledBase, scaledCoeffs, commonDen.Int64(), nil
}
