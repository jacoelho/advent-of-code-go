package aoc2025

import (
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

func parseBatteryBank(r io.Reader) ([][]int, error) {
	s := scanner.NewScanner(r, func(b []byte) ([]int, error) {
		digits := make([]int, len(b))
		for i, ch := range b {
			digits[i] = int(ch - '0')
		}
		return digits, nil
	})
	return slices.Collect(s.Values()), s.Err()
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
	banks, err := parseBatteryBank(r)
	if err != nil {
		return "", err
	}

	total := 0
	for _, bank := range banks {
		total += largestJoltage(bank, numBatteries)
	}
	return strconv.Itoa(total), nil
}

func day03p01(r io.Reader) (string, error) {
	return day03(r, 2)
}

func day03p02(r io.Reader) (string, error) {
	return day03(r, 12)
}
