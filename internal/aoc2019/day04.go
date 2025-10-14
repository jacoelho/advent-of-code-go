package aoc2019

import (
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/internal/convert"
	"github.com/jacoelho/advent-of-code-go/internal/scanner"
	"github.com/jacoelho/advent-of-code-go/internal/xmaps"
	"github.com/jacoelho/advent-of-code-go/internal/xslices"
)

func parse(r io.Reader) (int, int, error) {
	s := scanner.NewScanner(r, func(b []byte) (string, error) {
		return string(b), nil
	})

	var line string
	for l := range s.Values() {
		line = l
		break
	}

	if err := s.Err(); err != nil {
		return 0, 0, err
	}

	nums := convert.ExtractDigits[int](line)
	if len(nums) != 2 {
		return 0, 0, fmt.Errorf("expected 2 numbers, got %d", len(nums))
	}

	return nums[0], nums[1], nil
}

func hasPair(digits []int) bool {
	return xmaps.Any(func(k int, v int) bool { return v >= 2 }, xslices.Frequencies(digits))
}

func isNonDecreasing(digits []int) bool {
	for pair := range xslices.Window(2, digits) {
		if pair[0] > pair[1] {
			return false
		}
	}
	return true
}

func hasExactlyOnePair(digits []int) bool {
	return xmaps.Any(func(k int, v int) bool { return v == 2 }, xslices.Frequencies(digits))
}

func day4p01(r io.Reader) (string, error) {
	min, max, err := parse(r)
	if err != nil {
		return "", err
	}

	count := 0
	for n := min; n <= max; n++ {
		digits := convert.ToDigits(n)
		if hasPair(digits) && isNonDecreasing(digits) {
			count++
		}
	}

	return strconv.Itoa(count), nil
}

func day4p02(r io.Reader) (string, error) {
	min, max, err := parse(r)
	if err != nil {
		return "", err
	}

	count := 0
	for n := min; n <= max; n++ {
		digits := convert.ToDigits(n)
		if hasExactlyOnePair(digits) && isNonDecreasing(digits) {
			count++
		}
	}

	return strconv.Itoa(count), nil
}
