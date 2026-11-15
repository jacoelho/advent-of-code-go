package aoc2024

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

type calibrationEquation struct {
	target int64
	values []int64
}

func parseCalibrationEquations(r io.Reader) ([]calibrationEquation, error) {
	var equations []calibrationEquation
	s := bufio.NewScanner(r)
	for line := 1; s.Scan(); line++ {
		target, operands, ok := strings.Cut(s.Text(), ":")
		if !ok {
			return nil, fmt.Errorf("line %d: calibration equation needs a colon", line)
		}
		value, err := calibrationNumber(strings.TrimSpace(target))
		if err != nil {
			return nil, fmt.Errorf("line %d: calibration target: %w", line, err)
		}
		equation := calibrationEquation{target: value}
		for field := range strings.FieldsSeq(operands) {
			value, err := calibrationNumber(field)
			if err != nil {
				return nil, fmt.Errorf("line %d: calibration operand: %w", line, err)
			}
			equation.values = append(equation.values, value)
		}
		if len(equation.values) == 0 {
			return nil, fmt.Errorf("line %d: calibration equation needs an operand", line)
		}
		equations = append(equations, equation)
	}
	return equations, s.Err()
}

func calibrationNumber(field string) (int64, error) {
	value, err := convert.ScanNumber[int64]([]byte(field))
	if err != nil {
		return 0, err
	}
	if value < 0 {
		return 0, fmt.Errorf("%d is negative", value)
	}
	return value, nil
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

// calibrationReachable undoes operators from the last operand backwards, so
// intermediate values only shrink and cannot overflow.
func calibrationReachable(values []int64, target int64, concat bool) bool {
	last := len(values) - 1
	if last == 0 {
		return target == values[0]
	}
	operand, prefix := values[last], values[:last]
	if operand == 0 {
		// Multiplying by zero yields zero whatever the prefix evaluates to.
		if target == 0 {
			return true
		}
	} else if target%operand == 0 && calibrationReachable(prefix, target/operand, concat) {
		return true
	}
	if target >= operand && calibrationReachable(prefix, target-operand, concat) {
		return true
	}
	if concat {
		if rest, ok := removeCalibrationSuffix(target, operand); ok {
			return calibrationReachable(prefix, rest, concat)
		}
	}
	return false
}

func calibrationTotal(r io.Reader, concat bool) (string, error) {
	equations, err := parseCalibrationEquations(r)
	if err != nil {
		return "", err
	}
	var total int64
	for _, equation := range equations {
		if !calibrationReachable(equation.values, equation.target, concat) {
			continue
		}
		sum, ok := xmath.CheckedAdd(total, equation.target)
		if !ok {
			return "", fmt.Errorf("calibration total overflows int64")
		}
		total = sum
	}
	return strconv.FormatInt(total, 10), nil
}

func day07p01(r io.Reader) (string, error) {
	return calibrationTotal(r, false)
}

func day07p02(r io.Reader) (string, error) {
	return calibrationTotal(r, true)
}
