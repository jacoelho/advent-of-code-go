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

const noSolution = -1

type machineRequirements struct {
	lightsMask  uint
	buttonMasks []uint
	joltage     []int
}

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

func parseJoltageToken(token string) ([]int, error) {
	if len(token) < 2 || token[0] != '{' || token[len(token)-1] != '}' {
		return nil, fmt.Errorf("invalid joltage token: %s", token)
	}

	return convert.ExtractDigits[int](token[1 : len(token)-1]), nil
}

func parseInitializationProcedure(r io.Reader) ([]machineRequirements, error) {
	s := scanner.NewScanner(r, func(line []byte) (machineRequirements, error) {
		tokens := strings.Fields(string(line))

		lightsMask, err := parseLightsToken(tokens[0])
		if err != nil {
			return machineRequirements{}, err
		}

		joltage, err := parseJoltageToken(tokens[len(tokens)-1])
		if err != nil {
			return machineRequirements{}, err
		}

		buttonMasks, err := parseButtonTokens(tokens[1 : len(tokens)-1])
		if err != nil {
			return machineRequirements{}, err
		}

		return machineRequirements{
			lightsMask:  lightsMask,
			buttonMasks: buttonMasks,
			joltage:     joltage,
		}, nil
	})

	return slices.Collect(s.Values()), s.Err()
}

// validateAndSum validates button values and computes their sum.
// returns true if all values are non-negative and divisible by commonDen.
func validateAndSum(scaledBase []int64, scaledCoeffs [][]int64, params []int, commonDen int64) (int, bool) {
	total := 0
	for i := range scaledBase {
		value := scaledBase[i]
		for k := range params {
			value += scaledCoeffs[k][i] * int64(params[k])
		}
		if value < 0 || value%commonDen != 0 {
			return 0, false
		}
		total += int(value / commonDen)
	}
	return total, true
}

// buildCoefficientMatrix builds the coefficient matrix for the button press constraints.
func buildCoefficientMatrix(step machineRequirements) *matrix.Matrix {
	buttonCount := len(step.buttonMasks)
	counterCount := len(step.joltage)

	mat := matrix.New(counterCount, buttonCount+1)
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

// computeFreeVarLimits computes bounds for free variables.
// first computes upper bounds based on minimum joltage, then tightens them
// by considering constraints from negative coefficients.
func computeFreeVarLimits(step machineRequirements, scaledBase []int64, scaledCoeffs [][]int64, freeCols []int) []int {
	buttonCount := len(step.buttonMasks)
	freeCount := len(freeCols)

	// compute upper bounds for each button based on minimum joltage it affects
	upper := make([]int, buttonCount)
	for buttonIndex, mask := range step.buttonMasks {
		if mask == 0 {
			upper[buttonIndex] = 0
			continue
		}
		minJoltage := -1
		for counterIndex := 0; counterIndex < len(step.joltage); counterIndex++ {
			if mask&(1<<counterIndex) != 0 {
				if minJoltage == -1 || step.joltage[counterIndex] < minJoltage {
					minJoltage = step.joltage[counterIndex]
				}
			}
		}
		if minJoltage == -1 {
			upper[buttonIndex] = 0
		} else {
			upper[buttonIndex] = minJoltage
		}
	}

	// tighten bounds for free variables based on constraints
	limits := make([]int, freeCount)
	for freeVarIndex := range freeCount {
		limit := upper[freeCols[freeVarIndex]]
		for buttonIndex := range buttonCount {
			coeff := scaledCoeffs[freeVarIndex][buttonIndex]
			if coeff >= 0 {
				continue
			}

			// compute maximum help from other free variables
			maxHelp := int64(0)
			for otherFreeVarIndex := range freeCount {
				if otherFreeVarIndex == freeVarIndex {
					continue
				}
				otherCoeff := scaledCoeffs[otherFreeVarIndex][buttonIndex]
				if otherCoeff > 0 {
					maxHelp += otherCoeff * int64(upper[freeCols[otherFreeVarIndex]])
				}
			}

			maxAllowed := max((scaledBase[buttonIndex]+maxHelp)/-coeff, 0)
			if int64(limit) > maxAllowed {
				limit = int(maxAllowed)
			}
		}
		limits[freeVarIndex] = max(limit, 0)
	}

	return limits
}

// findMinimum performs a DFS search to find the minimum valid sum of button values.
// returns noSolution if no valid solution exists.
func findMinimum(scaledBase []int64, scaledCoeffs [][]int64, commonDen int64, limits []int) int {
	freeCount := len(limits)
	best := noSolution
	params := make([]int, freeCount)

	var dfs func(freeVarIndex int)
	dfs = func(freeVarIndex int) {
		if freeVarIndex == freeCount {
			total := 0
			for i := range scaledBase {
				value := scaledBase[i]
				for k := range params {
					value += scaledCoeffs[k][i] * int64(params[k])
				}
				if value < 0 || value%commonDen != 0 {
					return
				}
				total += int(value / commonDen)
			}
			if best == noSolution || total < best {
				best = total
			}
			return
		}

		limit := limits[freeVarIndex]
		for paramValue := 0; paramValue <= limit; paramValue++ {
			params[freeVarIndex] = paramValue

			feasible := true
			for buttonIndex := range scaledBase {
				value := scaledBase[buttonIndex]
				for fv := 0; fv <= freeVarIndex; fv++ {
					value += scaledCoeffs[fv][buttonIndex] * int64(params[fv])
				}
				for fv := freeVarIndex + 1; fv < freeCount; fv++ {
					coeff := scaledCoeffs[fv][buttonIndex]
					if coeff > 0 {
						value += coeff * int64(limits[fv])
					}
				}
				if value < 0 {
					feasible = false
					break
				}
			}
			if !feasible {
				continue
			}

			dfs(freeVarIndex + 1)
		}
	}

	dfs(0)
	return best
}

// minButtonPresses finds the minimum number of button presses needed
// to achieve the target lights configuration.
func minButtonPresses(step machineRequirements) int {
	target := step.lightsMask
	n := len(step.buttonMasks)

	minPresses := n + 1
	for buttonMask := uint(0); buttonMask < (1 << n); buttonMask++ {
		count := bits.OnesCount(buttonMask)
		if count >= minPresses {
			continue
		}

		var result uint
		remainingMask := buttonMask
		for remainingMask != 0 {
			buttonIndex := bits.TrailingZeros(remainingMask)
			result ^= step.buttonMasks[buttonIndex]
			remainingMask &= remainingMask - 1
		}

		if result == target {
			minPresses = count
		}
	}
	return minPresses
}

// solveSystem solves the system of equations for the button press constraints.
func solveSystem(step machineRequirements) (scaledBase []int64, scaledCoeffs [][]int64, commonDen int64, freeCols []int, err error) {
	buttonCount := len(step.buttonMasks)
	mat := buildCoefficientMatrix(step)

	result := matrix.RREF(mat)
	if result.Inconsistent {
		return nil, nil, 0, nil, fmt.Errorf("inconsistent system")
	}

	sol := matrix.ExtractParametricSolution(mat, result.PivotCols, buttonCount)
	scaledBase, scaledCoeffs, commonDen = matrix.ScaleToIntegers(sol.Base, sol.Coeffs)
	return scaledBase, scaledCoeffs, commonDen, sol.FreeCols, nil
}

// handleFreeVariables sets up the parametric system and finds the minimum solution when free variables exist.
func handleFreeVariables(step machineRequirements, scaledBase []int64, scaledCoeffs [][]int64, commonDen int64, freeCols []int) int {
	limits := computeFreeVarLimits(step, scaledBase, scaledCoeffs, freeCols)
	return findMinimum(scaledBase, scaledCoeffs, commonDen, limits)
}

// minButtonPressesForJoltage finds the minimum number of button presses
// needed to satisfy the joltage constraints.
func minButtonPressesForJoltage(step machineRequirements) int {
	buttonCount := len(step.buttonMasks)
	if buttonCount == 0 || len(step.joltage) == 0 {
		return 0
	}

	scaledBase, scaledCoeffs, commonDen, freeCols, err := solveSystem(step)
	if err != nil {
		return noSolution
	}

	freeCount := len(freeCols)
	if freeCount == 0 {
		total, valid := validateAndSum(scaledBase, scaledCoeffs, nil, commonDen)
		if !valid {
			return noSolution
		}
		return total
	}

	return handleFreeVariables(step, scaledBase, scaledCoeffs, commonDen, freeCols)
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
