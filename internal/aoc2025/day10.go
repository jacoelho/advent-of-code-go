package aoc2025

import (
	"fmt"
	"io"
	"math/bits"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/matrix"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xslices"
)

type initializationStep struct {
	lightsMask  uint
	buttonMasks []uint
	joltage     []int
}

// parseLightsToken parses a lights token like "[######]"
func parseLightsToken(token string) (mask uint, err error) {
	if len(token) < 2 || token[0] != '[' || token[len(token)-1] != ']' {
		return 0, fmt.Errorf("invalid lights token: %s", token)
	}

	pattern := token[1 : len(token)-1]
	width := 0
	for _, ch := range pattern {
		switch ch {
		case '#':
			mask |= 1 << width
			width++
		case '.':
			width++
		default:
			return 0, fmt.Errorf("invalid character '%c' in lights", ch)
		}
	}

	return mask, nil
}

// parseButtonTokens parses button tokens like "(0,1,3,5)"
func parseButtonTokens(tokens []string) ([]uint, error) {
	buttonMasks := make([]uint, 0, len(tokens))

	for _, token := range tokens {
		if len(token) < 2 || token[0] != '(' || token[len(token)-1] != ')' {
			return nil, fmt.Errorf("invalid button token: %s", token)
		}

		indices := convert.ExtractDigits[int](token[1 : len(token)-1])

		var mask uint
		for _, idx := range indices {
			mask |= 1 << idx
		}

		buttonMasks = append(buttonMasks, mask)
	}

	return buttonMasks, nil
}

// parseJoltageToken parses a joltage token like "{21,37,18,9,8,9}"
func parseJoltageToken(token string) ([]int, error) {
	if len(token) < 2 || token[0] != '{' || token[len(token)-1] != '}' {
		return nil, fmt.Errorf("invalid joltage token: %s", token)
	}

	return convert.ExtractDigits[int](token[1 : len(token)-1]), nil
}

func parseInitializationProcedure(r io.Reader) ([]initializationStep, error) {
	s := scanner.NewScanner(r, func(line []byte) (initializationStep, error) {
		tokens := strings.Fields(string(line))
		if len(tokens) < 2 {
			return initializationStep{}, fmt.Errorf("invalid input format")
		}

		// Parse first token - lights
		lightsMask, err := parseLightsToken(tokens[0])
		if err != nil {
			return initializationStep{}, err
		}

		// Parse last token - joltage
		joltage, err := parseJoltageToken(tokens[len(tokens)-1])
		if err != nil {
			return initializationStep{}, err
		}

		// Parse middle tokens - buttons
		buttonMasks, err := parseButtonTokens(tokens[1 : len(tokens)-1])
		if err != nil {
			return initializationStep{}, err
		}

		return initializationStep{
			lightsMask:  lightsMask,
			buttonMasks: buttonMasks,
			joltage:     joltage,
		}, nil
	})

	return slices.Collect(s.Values()), s.Err()
}

// computeButtonValue computes the value for button i given parameter assignments
func computeButtonValue(i int, scaledBase []int64, scaledCoeffs [][]int64, params []int) int64 {
	value := scaledBase[i]
	for k := range params {
		value += scaledCoeffs[k][i] * int64(params[k])
	}
	return value
}

// validateAndSum validates button values and computes their sum.
// Returns (sum, valid) where valid indicates all values are non-negative and divisible by commonDen.
func validateAndSum(scaledBase []int64, scaledCoeffs [][]int64, params []int, commonDen int64) (int, bool) {
	total := 0
	for i := range scaledBase {
		value := computeButtonValue(i, scaledBase, scaledCoeffs, params)
		if value < 0 || value%commonDen != 0 {
			return 0, false
		}
		total += int(value / commonDen)
	}
	return total, true
}

// computeMaxPotential computes maximum potential value for button i
// considering assigned parameters and upper bounds on unassigned parameters.
func computeMaxPotential(i int, scaledBase []int64, scaledCoeffs [][]int64, params []int, limits []int, assigned int) int64 {
	value := scaledBase[i]

	// Add contributions from assigned parameters
	for k := 0; k <= assigned; k++ {
		value += scaledCoeffs[k][i] * int64(params[k])
	}

	// Add maximum potential from unassigned parameters
	for k := assigned + 1; k < len(params); k++ {
		coeff := scaledCoeffs[k][i]
		if coeff > 0 {
			value += coeff * int64(limits[k])
		}
	}

	return value
}

func buildCoefficientMatrix(step initializationStep) *matrix.Matrix[matrix.Rat] {
	buttonCount := len(step.buttonMasks)
	counterCount := len(step.joltage)

	mat := matrix.New[matrix.Rat](counterCount, buttonCount+1)
	for i := range counterCount {
		mat.Set(i, buttonCount, matrix.NewRat(int64(step.joltage[i]), 1))
	}

	for j, mask := range step.buttonMasks {
		for idx := range counterCount {
			if mask&(1<<idx) != 0 {
				mat.Set(idx, j, matrix.NewRat(1, 1))
			}
		}
	}

	return mat
}

func computeUpperBounds(step initializationStep) []int {
	buttonCount := len(step.buttonMasks)
	upper := make([]int, buttonCount)

	for j, mask := range step.buttonMasks {
		if mask == 0 {
			upper[j] = 0
			continue
		}

		// Find minimum joltage affected by this button
		minJoltage := -1
		for idx := 0; idx < len(step.joltage); idx++ {
			if mask&(1<<idx) != 0 {
				if minJoltage == -1 || step.joltage[idx] < minJoltage {
					minJoltage = step.joltage[idx]
				}
			}
		}

		if minJoltage == -1 {
			upper[j] = 0
		} else {
			upper[j] = minJoltage
		}
	}

	return upper
}

func computeFreeVarLimits(scaledBase []int64, scaledCoeffs [][]int64, freeCols []int, upper []int) []int {
	buttonCount := len(scaledBase)
	freeCount := len(freeCols)
	limits := make([]int, freeCount)

	for idx := range freeCols {
		limit := upper[freeCols[idx]]
		for i := range buttonCount {
			c := scaledCoeffs[idx][i]
			if c >= 0 {
				continue
			}

			maxHelp := int64(0)
			for k := range freeCount {
				if k == idx {
					continue
				}
				coeff := scaledCoeffs[k][i]
				if coeff > 0 {
					maxHelp += coeff * int64(upper[freeCols[k]])
				}
			}

			maxAllowed := max((scaledBase[i]+maxHelp)/-c, 0)
			if int64(limit) > maxAllowed {
				limit = int(maxAllowed)
			}
		}
		limits[idx] = max(limit, 0)
	}

	return limits
}

func minButtonPresses(step initializationStep) int {
	target := step.lightsMask
	n := len(step.buttonMasks)

	minPresses := n + 1
	for mask := uint(0); mask < (1 << n); mask++ {
		count := bits.OnesCount(mask)
		if count >= minPresses {
			continue
		}

		var result uint
		m := mask
		for m != 0 {
			idx := bits.TrailingZeros(m)
			result ^= step.buttonMasks[idx]
			m &= m - 1
		}

		if result == target {
			minPresses = count
		}
	}
	return minPresses
}

func minButtonPressesForJoltage(step initializationStep) int {
	buttonCount := len(step.buttonMasks)
	if buttonCount == 0 || len(step.joltage) == 0 {
		return 0
	}

	mat := buildCoefficientMatrix(step)

	result := matrix.RREF(mat)
	if result.Inconsistent {
		return -1
	}

	sol := matrix.ExtractParametricSolution(mat, result.PivotCols, buttonCount)
	scaledBase, scaledCoeffs, commonDen := matrix.ScaleToIntegers(sol.Base, sol.Coeffs)

	freeCount := len(sol.FreeCols)

	// If there are no free variables, compute directly
	if freeCount == 0 {
		total, valid := validateAndSum(scaledBase, scaledCoeffs, nil, commonDen)
		if !valid {
			return -1
		}
		return total
	}

	upper := computeUpperBounds(step)
	limits := computeFreeVarLimits(scaledBase, scaledCoeffs, sol.FreeCols, upper)

	// DFS to find minimum button presses
	best := -1
	params := make([]int, freeCount)

	feasible := func(assigned int) bool {
		for i := range buttonCount {
			maxNum := computeMaxPotential(i, scaledBase, scaledCoeffs, params, limits, assigned)
			if maxNum < 0 {
				return false
			}
		}
		return true
	}

	var dfs func(idx int)
	dfs = func(idx int) {
		if idx == freeCount {
			total, valid := validateAndSum(scaledBase, scaledCoeffs, params, commonDen)
			if !valid {
				return
			}
			if best == -1 || total < best {
				best = total
			}
			return
		}

		for v := 0; v <= limits[idx]; v++ {
			params[idx] = v
			if !feasible(idx) {
				continue
			}
			dfs(idx + 1)
		}
	}

	dfs(0)
	return best
}

func day10p01(r io.Reader) (string, error) {
	steps, err := parseInitializationProcedure(r)
	if err != nil {
		return "", err
	}

	total := xslices.Sum(xslices.Map(minButtonPresses, steps))
	return strconv.Itoa(total), nil
}

func day10p02(r io.Reader) (string, error) {
	steps, err := parseInitializationProcedure(r)
	if err != nil {
		return "", err
	}

	total := xslices.Sum(xslices.Map(minButtonPressesForJoltage, steps))
	return strconv.Itoa(total), nil
}
