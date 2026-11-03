package aoc2025

import (
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
)

func parseBatteryBank(b []byte, numBatteries int) ([]int, error) {
	if len(b) < numBatteries {
		return nil, fmt.Errorf("battery bank needs at least %d digits", numBatteries)
	}
	digits := make([]int, len(b))
	for i, ch := range b {
		if ch < '0' || ch > '9' {
			return nil, fmt.Errorf("invalid battery digit %q", ch)
		}
		digits[i] = int(ch - '0')
	}
	return digits, nil
}

func largestJoltage(bank []int, k int) int {
	toRemove := len(bank) - k
	stack := make([]int, 0, k)

	for _, d := range bank {
		for len(stack) > 0 && toRemove > 0 && stack[len(stack)-1] < d {
			stack = stack[:len(stack)-1]
			toRemove--
		}
		stack = append(stack, d)
	}

	return convert.FromDigits(stack[:k])
}

func day03(r io.Reader, numBatteries int) (string, error) {
	s := scanner.NewScanner(r, func(b []byte) ([]int, error) {
		return parseBatteryBank(b, numBatteries)
	})
	total := xiter.From(s.Values()).SumBy(func(bank []int) int {
		return largestJoltage(bank, numBatteries)
	})
	if err := s.Err(); err != nil {
		return "", err
	}
	return strconv.Itoa(total), nil
}

func day03p01(r io.Reader) (string, error) {
	return day03(r, 2)
}

func day03p02(r io.Reader) (string, error) {
	return day03(r, 12)
}
