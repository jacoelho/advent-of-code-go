package aoc2022

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

func calorieTotals(r io.Reader) ([3]int64, error) {
	var top [3]int64
	input, err := puzzleInput(r)
	if err != nil {
		return top, err
	}
	for group := range strings.SplitSeq(input, "\n\n") {
		var total int64
		for line := range strings.SplitSeq(group, "\n") {
			value, err := convert.ScanNumber[int64]([]byte(line))
			if err != nil {
				return top, fmt.Errorf("calorie line %q: %w", line, err)
			}
			sum, ok := xmath.CheckedAdd(total, value)
			if value < 0 || !ok {
				return top, fmt.Errorf("invalid calorie total")
			}
			total = sum
		}
		for i := range top {
			if total > top[i] {
				total, top[i] = top[i], total
			}
		}
	}
	return top, nil
}

func day01p01(r io.Reader) (string, error) {
	top, err := calorieTotals(r)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(top[0], 10), nil
}

func day01p02(r io.Reader) (string, error) {
	top, err := calorieTotals(r)
	if err != nil {
		return "", err
	}
	var total int64
	for _, value := range top {
		sum, ok := xmath.CheckedAdd(total, value)
		if !ok {
			return "", fmt.Errorf("calorie sum overflows")
		}
		total = sum
	}
	return strconv.FormatInt(total, 10), nil
}
