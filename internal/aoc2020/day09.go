package aoc2020

import (
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

func parseXMAS(r io.Reader) ([]int64, error) {
	s := scanner.NewScanner(r, func(line []byte) (int64, error) {
		number, err := strconv.ParseInt(string(line), 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid XMAS number %q: %w", line, err)
		}
		return number, nil
	})
	return slices.Collect(s.Values()), s.Err()
}

func invalidXMASNumber(numbers []int64, preamble int) (int64, error) {
	if preamble < 2 || len(numbers) <= preamble {
		return 0, fmt.Errorf("XMAS data requires at least %d preamble numbers and a candidate", preamble)
	}
	for i := preamble; i < len(numbers); i++ {
		seen := make(map[int64]bool, preamble)
		valid := false
		for _, number := range numbers[i-preamble : i] {
			if complement := numbers[i] - number; complement != number && seen[complement] {
				valid = true
				break
			}
			seen[number] = true
		}
		if !valid {
			return numbers[i], nil
		}
	}
	return 0, fmt.Errorf("all XMAS numbers follow the preamble rule")
}

func xmasWeakness(numbers []int64, target int64) (int64, error) {
	// Keep the earliest prefix index so a repeated sum cannot hide a valid range.
	firstPrefix := map[int64]int{0: 0}
	var sum int64
	for i, number := range numbers {
		sum += number
		if start, ok := firstPrefix[sum-target]; ok && i+1-start >= 2 {
			window := numbers[start : i+1]
			return slices.Min(window) + slices.Max(window), nil
		}
		if _, exists := firstPrefix[sum]; !exists {
			firstPrefix[sum] = i + 1
		}
	}
	return 0, fmt.Errorf("no contiguous XMAS range of at least two numbers sums to %d", target)
}

func day09p01(r io.Reader) (string, error) {
	numbers, err := parseXMAS(r)
	if err != nil {
		return "", err
	}
	invalid, err := invalidXMASNumber(numbers, 25)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(invalid, 10), nil
}

func day09p02(r io.Reader) (string, error) {
	numbers, err := parseXMAS(r)
	if err != nil {
		return "", err
	}
	invalid, err := invalidXMASNumber(numbers, 25)
	if err != nil {
		return "", err
	}
	weakness, err := xmasWeakness(numbers, invalid)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(weakness, 10), nil
}
