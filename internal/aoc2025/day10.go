package aoc2025

import (
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"math/bits"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/matrix"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

const (
	noSolution        = -1
	maxMachineWidth   = 64
	maxMachineButtons = 24
	maxMachines       = 1024
	maxMachineWork    = 250_000_000
)

type machineRequirements struct {
	lightsMask  uint64
	buttonMasks []uint64
	joltage     []int64
}

func parseLightsToken(token string) (uint64, error) {
	if len(token) < 3 || token[0] != '[' || token[len(token)-1] != ']' {
		return 0, fmt.Errorf("invalid lights token %q", token)
	}
	pattern := token[1 : len(token)-1]
	if len(pattern) > maxMachineWidth {
		return 0, fmt.Errorf("machine exceeds %d lights", maxMachineWidth)
	}

	var mask uint64
	for i, ch := range pattern {
		switch ch {
		case '#':
			mask |= uint64(1) << i
		case '.':
		default:
			return 0, fmt.Errorf("invalid light %q", ch)
		}
	}
	return mask, nil
}

func parseButtonTokens(tokens []string, width int) ([]uint64, error) {
	if len(tokens) > maxMachineButtons {
		return nil, fmt.Errorf("machine exceeds %d buttons", maxMachineButtons)
	}
	buttonMasks := make([]uint64, 0, len(tokens))
	for _, token := range tokens {
		if len(token) < 2 || token[0] != '(' || token[len(token)-1] != ')' {
			return nil, fmt.Errorf("invalid button token %q", token)
		}
		var mask uint64
		if token != "()" {
			for field := range strings.SplitSeq(token[1:len(token)-1], ",") {
				idx, err := convert.ScanNumber[int]([]byte(field))
				if err != nil {
					return nil, fmt.Errorf("button index: %w", err)
				}
				if idx < 0 || idx >= width {
					return nil, fmt.Errorf("button index %d outside %d lights", idx, width)
				}
				bit := uint64(1) << idx
				if mask&bit != 0 {
					return nil, fmt.Errorf("duplicate button index %d", idx)
				}
				mask |= bit
			}
		}
		buttonMasks = append(buttonMasks, mask)
	}
	return buttonMasks, nil
}

func parseJoltageToken(token string, width int) ([]int64, error) {
	if len(token) < 3 || token[0] != '{' || token[len(token)-1] != '}' {
		return nil, fmt.Errorf("invalid joltage token %q", token)
	}
	var joltage []int64
	for field := range strings.SplitSeq(token[1:len(token)-1], ",") {
		value, err := convert.ScanNumber[int64]([]byte(field))
		if err != nil {
			return nil, fmt.Errorf("joltage count: %w", err)
		}
		if value < 0 {
			return nil, fmt.Errorf("joltage count must be nonnegative")
		}
		joltage = append(joltage, value)
		if len(joltage) > width {
			return nil, fmt.Errorf("joltage width differs from lights")
		}
	}
	if len(joltage) != width {
		return nil, fmt.Errorf("joltage width differs from lights")
	}
	return joltage, nil
}

func parseMachine(line []byte) (machineRequirements, error) {
	tokens := strings.Fields(string(line))
	if len(tokens) < 2 {
		return machineRequirements{}, fmt.Errorf("machine needs lights and joltage")
	}
	lights, err := parseLightsToken(tokens[0])
	if err != nil {
		return machineRequirements{}, err
	}
	width := len(tokens[0]) - 2
	joltage, err := parseJoltageToken(tokens[len(tokens)-1], width)
	if err != nil {
		return machineRequirements{}, err
	}
	buttons, err := parseButtonTokens(tokens[1:len(tokens)-1], width)
	if err != nil {
		return machineRequirements{}, err
	}
	return machineRequirements{lightsMask: lights, buttonMasks: buttons, joltage: joltage}, nil
}

func parseInitializationProcedure(r io.Reader) ([]machineRequirements, error) {
	machineIndex := 0
	s := scanner.NewScanner(r, func(line []byte) (machineRequirements, error) {
		machineIndex++
		if machineIndex > maxMachines {
			return machineRequirements{}, fmt.Errorf("machine %d: procedure exceeds %d machines", machineIndex, maxMachines)
		}
		machine, err := parseMachine(line)
		if err != nil {
			return machineRequirements{}, fmt.Errorf("machine %d: %w", machineIndex, err)
		}
		return machine, nil
	})
	steps := slices.Collect(s.Values())
	return steps, s.Err()
}

type joltageFamily struct {
	base         []int64
	coefficients [][]int64
	denominator  int64
	freeButtons  []int
}

func (f *joltageFamily) validateAndSum(params []int64) (int64, bool) {
	var total int64
	for i, value := range f.base {
		for k, param := range params {
			value += f.coefficients[k][i] * param
		}
		if value < 0 || value%f.denominator != 0 {
			return 0, false
		}
		total += value / f.denominator
	}
	return total, true
}

func (f *joltageFamily) freeVarLimits(step machineRequirements) []int64 {
	upper := make([]int64, len(step.buttonMasks))
	for i, mask := range step.buttonMasks {
		if mask == 0 {
			continue
		}
		upper[i] = math.MaxInt64
		for counter, target := range step.joltage {
			if mask&(uint64(1)<<counter) != 0 {
				upper[i] = min(upper[i], target)
			}
		}
	}

	limits := make([]int64, len(f.freeButtons))
	for k, col := range f.freeButtons {
		limit := upper[col]
		for i, value := range f.base {
			coefficient := f.coefficients[k][i]
			if coefficient >= 0 {
				continue
			}
			// Maximize the other positive terms to bound this free variable while
			// keeping the dependent variable nonnegative.
			available := big.NewInt(value)
			for other, freeCol := range f.freeButtons {
				if c := f.coefficients[other][i]; other != k && c > 0 {
					term := new(big.Int).Mul(big.NewInt(c), big.NewInt(upper[freeCol]))
					available.Add(available, term)
				}
			}
			divisor := new(big.Int).Neg(big.NewInt(coefficient))
			allowed := new(big.Int).Quo(available, divisor)
			if allowed.Sign() < 0 {
				limit = 0
			} else if allowed.Cmp(big.NewInt(limit)) < 0 {
				limit = allowed.Int64()
			}
		}
		limits[k] = limit
	}
	return limits
}

func (f *joltageFamily) admitArithmetic(limits []int64) error {
	totalBound := new(big.Int)
	for i, value := range f.base {
		bound := new(big.Int).Abs(big.NewInt(value))
		for k, limit := range limits {
			term := new(big.Int).Abs(big.NewInt(f.coefficients[k][i]))
			term.Mul(term, big.NewInt(limit))
			bound.Add(bound, term)
		}
		// This absolute bound also covers every partial sum used by pruning.
		if !bound.IsInt64() {
			return fmt.Errorf("unsupported joltage range: scaled expression exceeds int64")
		}
		totalBound.Add(totalBound, bound.Quo(bound, big.NewInt(f.denominator)))
	}
	if !totalBound.IsInt64() {
		return fmt.Errorf("unsupported joltage range: press sum exceeds int64")
	}
	return nil
}

func spendMachineWork(remaining *int, cost int) error {
	if cost > *remaining {
		return fmt.Errorf("machine search work limit exceeded")
	}
	*remaining -= cost
	return nil
}

func (f *joltageFamily) minimum(limits []int64, remaining *int) (int64, error) {
	best := int64(noSolution)
	params := make([]int64, len(limits))
	work := 1 + len(f.base)*(len(limits)+1)

	var dfs func(int) error
	dfs = func(k int) error {
		if k == len(limits) {
			if err := spendMachineWork(remaining, work); err != nil {
				return err
			}
			if total, valid := f.validateAndSum(params); valid && (best == noSolution || total < best) {
				best = total
			}
			return nil
		}

		for param := int64(0); ; param++ {
			if err := spendMachineWork(remaining, work); err != nil {
				return err
			}
			params[k] = param
			feasible := true
			for i, value := range f.base {
				for j := 0; j <= k; j++ {
					value += f.coefficients[j][i] * params[j]
				}
				for j := k + 1; j < len(limits); j++ {
					if f.coefficients[j][i] > 0 {
						value += f.coefficients[j][i] * limits[j]
					}
				}
				if value < 0 {
					feasible = false
					break
				}
			}
			if feasible {
				if err := dfs(k + 1); err != nil {
					return err
				}
			}
			if param == limits[k] {
				break
			}
		}
		return nil
	}
	if err := dfs(0); err != nil {
		return 0, err
	}
	return best, nil
}

func minButtonPresses(step machineRequirements, remaining *int) (int64, error) {
	n := len(step.buttonMasks)
	best := n + 1
	for subset := uint64(0); subset < uint64(1)<<n; subset++ {
		if err := spendMachineWork(remaining, n+1); err != nil {
			return 0, err
		}
		count := bits.OnesCount64(subset)
		if count >= best {
			continue
		}
		var result uint64
		for selected := subset; selected != 0; selected &= selected - 1 {
			result ^= step.buttonMasks[bits.TrailingZeros64(selected)]
		}
		if result == step.lightsMask {
			best = count
		}
	}
	if best > n {
		return 0, fmt.Errorf("lights have no reachable configuration")
	}
	return int64(best), nil
}

func solveJoltage(step machineRequirements) (joltageFamily, error) {
	buttons := make([]matrix.Var[int], len(step.buttonMasks))
	for i := range buttons {
		buttons[i] = matrix.Variable(i)
	}
	system := matrix.NewSystem(buttons...)
	for counter, target := range step.joltage {
		var sum matrix.Expr[int]
		for i, mask := range step.buttonMasks {
			if mask&(uint64(1)<<counter) != 0 {
				sum = sum.Add(buttons[i].Scale(1))
			}
		}
		system.Add(sum.Equal(target))
	}
	solution, err := system.Solve()
	if err != nil {
		if errors.Is(err, matrix.ErrInconsistent) {
			return joltageFamily{}, fmt.Errorf("joltage has inconsistent constraints: %w", err)
		}
		return joltageFamily{}, fmt.Errorf("joltage equations: %w", err)
	}

	free := solution.Free()
	family := joltageFamily{freeButtons: make([]int, len(free))}
	base := make([]matrix.Rat, len(buttons))
	coefficients := make([][]matrix.Rat, len(free))
	for i, button := range free {
		family.freeButtons[i] = button.Key()
		coefficients[i] = make([]matrix.Rat, len(buttons))
	}
	for i, button := range buttons {
		expression, err := solution.Expression(button)
		if err != nil {
			return joltageFamily{}, fmt.Errorf("joltage button %d: %w", i, err)
		}
		base[i] = expression.Constant()
		for j, parameter := range free {
			coefficients[j][i] = expression.Coefficient(parameter)
		}
	}
	family.base, family.coefficients, family.denominator, err = matrix.ScaleToIntegers(base, coefficients)
	if err != nil {
		return joltageFamily{}, fmt.Errorf("unsupported joltage range: %w", err)
	}
	return family, nil
}

func minButtonPressesForJoltage(step machineRequirements, remaining *int) (int64, error) {
	n := len(step.buttonMasks)
	if err := spendMachineWork(remaining, 1+len(step.joltage)*n*(n+1)); err != nil {
		return 0, err
	}

	family, err := solveJoltage(step)
	if err != nil {
		return 0, err
	}

	limits := family.freeVarLimits(step)
	if err := family.admitArithmetic(limits); err != nil {
		return 0, err
	}

	minimum, err := family.minimum(limits, remaining)
	if err != nil {
		return 0, err
	}
	if minimum == noSolution {
		return 0, fmt.Errorf("joltage has no nonnegative integer solution")
	}
	return minimum, nil
}

func day10p01(r io.Reader) (string, error) {
	steps, err := parseInitializationProcedure(r)
	if err != nil {
		return "", err
	}

	remaining := maxMachineWork
	var total int64
	for i, step := range steps {
		presses, err := minButtonPresses(step, &remaining)
		if err != nil {
			return "", fmt.Errorf("machine %d: %w", i+1, err)
		}
		total += presses
	}
	return strconv.FormatInt(total, 10), nil
}

func day10p02(r io.Reader) (string, error) {
	steps, err := parseInitializationProcedure(r)
	if err != nil {
		return "", err
	}

	remaining := maxMachineWork
	var total int64
	for i, step := range steps {
		presses, err := minButtonPressesForJoltage(step, &remaining)
		if err != nil {
			return "", fmt.Errorf("machine %d: %w", i+1, err)
		}
		if presses > math.MaxInt64-total {
			return "", fmt.Errorf("machine %d: total presses exceed int64", i+1)
		}
		total += presses
	}
	return strconv.FormatInt(total, 10), nil
}
