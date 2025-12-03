package aoc2025

import (
	"fmt"
	"io"
	"iter"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

type productRange struct {
	start int
	end   int
}

func (p productRange) Iter() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := p.start; i <= p.end; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

func parseProductRanges(r io.Reader) ([]productRange, error) {
	s := scanner.NewScannerWithSplit(r, scanner.SplitBySeparator([]byte(",")), func(b []byte) (productRange, error) {
		nums := convert.ExtractDigits[int](string(b))
		if len(nums) != 2 {
			return productRange{}, fmt.Errorf("invalid range: %s", b)
		}
		return productRange{start: nums[0], end: nums[1]}, nil
	})
	return slices.Collect(s.Values()), s.Err()
}

func day02(r io.Reader, filter func(int) bool) (string, error) {
	productRanges, err := parseProductRanges(r)
	if err != nil {
		return "", err
	}

	var sum int
	for _, pr := range productRanges {
		for id := range pr.Iter() {
			if filter(id) {
				sum += id
			}
		}
	}
	return strconv.Itoa(sum), nil
}

// hasRepeatedHalves returns true if the first half of the digits equals the second half.
func hasRepeatedHalves(n int) bool {
	digits := convert.ToDigits(n)
	d := len(digits)
	if d%2 != 0 {
		return false
	}
	half := d / 2
	return slices.Equal(digits[:half], digits[half:])
}

func day02p01(r io.Reader) (string, error) {
	return day02(r, hasRepeatedHalves)
}

// hasRepeatingPattern reports whether n's digits form a repeating pattern.
// A pattern repeats if the digit sequence can be split into identical blocks.
func hasRepeatingPattern(n int) bool {
	digits := convert.ToDigits(n)
	d := len(digits)
	if d < 2 {
		return false
	}

outer:
	for k := d / 2; k >= 1; k-- {
		if d%k != 0 {
			continue
		}
		for i := k; i < d; i++ {
			if digits[i] != digits[i%k] {
				continue outer
			}
		}
		return true
	}
	return false
}

func day02p02(r io.Reader) (string, error) {
	return day02(r, hasRepeatingPattern)
}
