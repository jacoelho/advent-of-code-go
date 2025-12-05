package aoc2025

import (
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

func parseInventory(r io.Reader) ([][2]int, []int, error) {
	s := scanner.NewScannerWithSplit(r, scanner.SplitBySeparator([]byte{'\n', '\n'}), func(b []byte) (string, error) {
		return string(b), nil
	})

	sections := slices.Collect(s.Values())
	if err := s.Err(); err != nil {
		return nil, nil, err
	}
	if len(sections) < 2 {
		return nil, nil, nil
	}

	var ranges [][2]int
	for line := range strings.SplitSeq(sections[0], "\n") {
		if line == "" {
			continue
		}
		nums := convert.ExtractDigits[int](line)
		if len(nums) != 2 {
			continue
		}
		ranges = append(ranges, [2]int{nums[0], nums[1]})
	}

	var integers []int
	for line := range strings.SplitSeq(sections[1], "\n") {
		if line == "" {
			continue
		}
		nums := convert.ExtractDigits[int](line)
		if len(nums) > 0 {
			integers = append(integers, nums[0])
		}
	}

	return ranges, integers, nil
}

func mergeRanges(ranges [][2]int) [][2]int {
	if len(ranges) == 0 {
		return nil
	}

	slices.SortFunc(ranges, func(a, b [2]int) int {
		return a[0] - b[0]
	})

	var merged [][2]int
	current := ranges[0]

	for i := 1; i < len(ranges); i++ {
		next := ranges[i]
		if next[0] <= current[1]+1 {
			current[1] = max(current[1], next[1])
		} else {
			merged = append(merged, current)
			current = next
		}
	}
	merged = append(merged, current)

	return merged
}

func inRange(ranges [][2]int, n int) bool {
	_, found := slices.BinarySearchFunc(ranges, n, func(r [2]int, target int) int {
		if target < r[0] {
			return 1
		}
		if target > r[1] {
			return -1
		}
		return 0
	})
	return found
}

func day05p01(r io.Reader) (string, error) {
	ranges, integers, err := parseInventory(r)
	if err != nil {
		return "", err
	}

	merged := mergeRanges(ranges)
	count := 0
	for _, n := range integers {
		if inRange(merged, n) {
			count++
		}
	}

	return strconv.Itoa(count), nil
}

func rangeLength(r [2]int) int {
	return r[1] - r[0] + 1
}

func day05p02(r io.Reader) (string, error) {
	ranges, _, err := parseInventory(r)
	if err != nil {
		return "", err
	}

	merged := mergeRanges(ranges)
	total := 0
	for _, r := range merged {
		total += rangeLength(r)
	}

	return strconv.Itoa(total), nil
}
