package aoc2024

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

type calibrationEquation struct {
	target int64
	values []int64
}

type calibrationOperators uint8

const (
	calibrationBasic calibrationOperators = iota
	calibrationWithConcat
)

type calibrationState struct {
	index int
	value int64
}

func parseCalibrationEquations(r io.Reader) ([]calibrationEquation, error) {
	var equations []calibrationEquation
	s := bufio.NewScanner(r)
	for s.Scan() {
		if len(equations) == 10_000 {
			return nil, fmt.Errorf("calibration exceeds 10000 equations")
		}
		target, operands, ok := strings.Cut(s.Text(), ":")
		if !ok {
			return nil, fmt.Errorf("calibration equation needs a colon")
		}
		value, err := convert.ScanNumber[int64]([]byte(strings.TrimSpace(target)))
		if err != nil {
			return nil, fmt.Errorf("calibration target: %w", err)
		}
		if value < 0 {
			return nil, fmt.Errorf("calibration target must be nonnegative")
		}
		equation := calibrationEquation{target: value}
		for field := range strings.FieldsSeq(operands) {
			value, err := convert.ScanNumber[int64]([]byte(field))
			if err != nil {
				return nil, fmt.Errorf("calibration operand: %w", err)
			}
			if value < 0 {
				return nil, fmt.Errorf("calibration operands must be nonnegative")
			}
			if len(equation.values) == 64 {
				return nil, fmt.Errorf("calibration equation exceeds 64 operands")
			}
			equation.values = append(equation.values, value)
		}
		if len(equation.values) == 0 {
			return nil, fmt.Errorf("calibration equation needs an operand")
		}
		equations = append(equations, equation)
	}
	return equations, s.Err()
}

func removeCalibrationSuffix(value, suffix int64) (int64, bool) {
	for {
		if value%10 != suffix%10 {
			return 0, false
		}
		value /= 10
		suffix /= 10
		if suffix == 0 {
			return value, true
		}
	}
}

func calibrationReachable(equation calibrationEquation, operators calibrationOperators, remaining *int) (bool, error) {
	failed := make(map[calibrationState]bool)
	var visit func(int, int64) (bool, error)
	visit = func(index int, value int64) (bool, error) {
		state := calibrationState{index: index, value: value}
		if failed[state] {
			return false, nil
		}
		if *remaining == 0 {
			return false, fmt.Errorf("calibration exceeds two million search states")
		}
		*remaining -= 1
		if index == 0 {
			return value == equation.values[0], nil
		}

		operand := equation.values[index]
		var candidates []int64
		if operand == 0 {
			// Any mathematical prefix, even above int64, can be reset by multiplying by zero.
			if value == 0 {
				return true, nil
			}
		} else if value%operand == 0 {
			candidates = append(candidates, value/operand)
		}
		if value >= operand {
			candidates = append(candidates, value-operand)
		}
		if operators == calibrationWithConcat {
			if prefix, ok := removeCalibrationSuffix(value, operand); ok {
				candidates = append(candidates, prefix)
			}
		}
		for _, candidate := range candidates {
			found, err := visit(index-1, candidate)
			if err != nil || found {
				return found, err
			}
		}
		failed[state] = true
		return false, nil
	}
	return visit(len(equation.values)-1, equation.target)
}

func calibrationTotal(r io.Reader, operators calibrationOperators) (string, error) {
	equations, err := parseCalibrationEquations(r)
	if err != nil {
		return "", err
	}
	remaining := 2_000_000
	var total int64
	for _, equation := range equations {
		found, err := calibrationReachable(equation, operators, &remaining)
		if err != nil {
			return "", err
		}
		if found {
			if equation.target > math.MaxInt64-total {
				return "", fmt.Errorf("calibration total overflows int64")
			}
			total += equation.target
		}
	}
	return strconv.FormatInt(total, 10), nil
}

func day07p01(r io.Reader) (string, error) {
	return calibrationTotal(r, calibrationBasic)
}

func day07p02(r io.Reader) (string, error) {
	return calibrationTotal(r, calibrationWithConcat)
}
