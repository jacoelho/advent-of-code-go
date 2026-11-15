package aoc2025

import (
	"cmp"
	"fmt"
	"io"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
)

func parseInventory(r io.Reader) ([][2]int, []int, error) {
	sections, err := scanner.ReadBlocks(r)
	if err != nil {
		return nil, nil, err
	}
	if len(sections) != 2 || len(sections[0]) == 0 {
		return nil, nil, fmt.Errorf("inventory must be ranges, a blank line, then ingredient IDs")
	}
	ranges := make([][2]int, len(sections[0]))
	for i, line := range sections[0] {
		low, high, ok := strings.Cut(line, "-")
		start, startErr := strconv.Atoi(low)
		end, endErr := strconv.Atoi(high)
		if !ok || startErr != nil || endErr != nil || start < 0 || start > end || end == math.MaxInt {
			return nil, nil, fmt.Errorf("invalid ingredient range %q", line)
		}
		ranges[i] = [2]int{start, end}
	}
	ids := make([]int, len(sections[1]))
	for i, line := range sections[1] {
		id, err := strconv.Atoi(line)
		if err != nil || id < 0 {
			return nil, nil, fmt.Errorf("invalid ingredient ID %q", line)
		}
		ids[i] = id
	}
	return ranges, ids, nil
}

func mergeRanges(ranges [][2]int) [][2]int {
	if len(ranges) == 0 {
		return nil
	}

	slices.SortFunc(ranges, func(a, b [2]int) int {
		return cmp.Compare(a[0], b[0])
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
	count := xiter.FromSlice(integers).Filter(func(n int) bool {
		return inRange(merged, n)
	}).Count()

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
	total := xiter.FromSlice(merged).SumBy(rangeLength)

	return strconv.Itoa(total), nil
}
