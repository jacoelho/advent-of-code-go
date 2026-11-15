package aoc2025

import (
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

const maxProductIDDigits = 18

type productRange struct {
	start int
	end   int
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

// repeatedBlockSum sums the IDs in p whose digits are one block repeated a
// number of times accepted by repeats. An ID of d digits made of a k-digit
// block h is h × (10^d − 1) / (10^k − 1), so the matching blocks for each (d, k)
// form a range. Numbers such as 1111 arise from several block lengths and are
// counted once.
func repeatedBlockSum(p productRange, repeats func(int) bool) (int, error) {
	if p.end >= pow10(maxProductIDDigits) {
		return 0, fmt.Errorf("product ID exceeds %d digits", maxProductIDDigits)
	}
	ids := make(map[int]struct{})
	for digits := 2; digits <= maxProductIDDigits; digits++ {
		low, high := max(p.start, pow10(digits-1)), min(p.end, pow10(digits)-1)
		if low > high {
			continue
		}
		for block := 1; block <= digits/2; block++ {
			if digits%block != 0 || !repeats(digits/block) {
				continue
			}
			multiplier := (pow10(digits) - 1) / (pow10(block) - 1)
			first := max(pow10(block-1), (low+multiplier-1)/multiplier)
			last := min(pow10(block)-1, high/multiplier)
			for h := first; h <= last; h++ {
				ids[h*multiplier] = struct{}{}
			}
		}
	}
	sum := 0
	for id := range ids {
		sum += id
	}
	return sum, nil
}

func pow10(n int) int {
	result := 1
	for range n {
		result *= 10
	}
	return result
}

func day02(r io.Reader, repeats func(int) bool) (string, error) {
	productRanges, err := parseProductRanges(r)
	if err != nil {
		return "", err
	}
	var sum int
	for _, pr := range productRanges {
		rangeSum, err := repeatedBlockSum(pr, repeats)
		if err != nil {
			return "", err
		}
		sum += rangeSum
	}
	return strconv.Itoa(sum), nil
}

func day02p01(r io.Reader) (string, error) {
	return day02(r, func(times int) bool { return times == 2 })
}

func day02p02(r io.Reader) (string, error) {
	return day02(r, func(times int) bool { return times >= 2 })
}
