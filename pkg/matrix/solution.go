package matrix

import "fmt"

// ScaleToIntegers multiplies base and coeffs by one positive common
// denominator. It returns the scaled base, scaled coefficients and denominator,
// or an error if any scaled value or the denominator cannot fit in int64.
func ScaleToIntegers(base []Rat, coeffs [][]Rat) ([]int64, [][]int64, int64, error) {
	commonDen := int64(1)
	includeDenominator := func(value Rat) bool {
		var denominator int64
		if value.value == nil {
			_, denominator = value.reducedCompact()
		} else {
			if !value.value.Denom().IsInt64() {
				return false
			}
			denominator = value.value.Denom().Int64()
		}
		if denominator == 1 {
			return true
		}
		divisor := int64(gcdMagnitudes(uint64(commonDen), uint64(denominator)))
		next, ok := multiplyIntegers(commonDen/divisor, denominator)
		if !ok {
			return false
		}
		commonDen = next
		return true
	}
	for _, value := range base {
		if !includeDenominator(value) {
			return nil, nil, 0, fmt.Errorf("common denominator exceeds int64")
		}
	}
	for _, row := range coeffs {
		for _, value := range row {
			if !includeDenominator(value) {
				return nil, nil, 0, fmt.Errorf("common denominator exceeds int64")
			}
		}
	}

	scale := func(values []Rat) ([]int64, error) {
		result := make([]int64, len(values))
		for i, value := range values {
			// Denominators already fit. A canonical big value therefore has an
			// oversized numerator, which a positive integer factor cannot shrink.
			if value.value != nil {
				return nil, fmt.Errorf("scaled value %d exceeds int64", i)
			}
			numerator, denominator := value.reducedCompact()
			integer, ok := multiplyIntegers(numerator, commonDen/denominator)
			if !ok {
				return nil, fmt.Errorf("scaled value %d exceeds int64", i)
			}
			result[i] = integer
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
	return scaledBase, scaledCoeffs, commonDen, nil
}
