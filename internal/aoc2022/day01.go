package aoc2022

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

func calorieTotals(r io.Reader) ([3]int64, error) {
	var top [3]int64
	input, err := puzzleInput(r)
	if err != nil {
		return top, err
	}
	var total int64
	finish := func() {
		for i := range top {
			if total > top[i] {
				total, top[i] = top[i], total
			}
		}
		total = 0
	}
	for line := range strings.SplitSeq(input, "\n") {
		if line == "" {
			finish()
			continue
		}
		value, err := convert.ScanNumber[int64]([]byte(line))
		if err != nil {
			return top, err
		}
		if value < 0 || total > math.MaxInt64-value {
			return top, fmt.Errorf("invalid calorie total")
		}
		total += value
	}
	finish()
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
		if total > math.MaxInt64-value {
			return "", fmt.Errorf("calorie sum overflows")
		}
		total += value
	}
	return strconv.FormatInt(total, 10), nil
}
