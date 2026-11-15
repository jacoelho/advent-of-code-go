package aoc2024

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

type pageOrderRules map[[2]int]bool

// compare relies on the puzzle guarantee that the rules order every pair of
// pages within an update; pages without a rule compare as equal.
func (rules pageOrderRules) compare(a, b int) int {
	switch {
	case rules[[2]int{a, b}]:
		return -1
	case rules[[2]int{b, a}]:
		return 1
	default:
		return 0
	}
}

func parsePageOrder(r io.Reader) (pageOrderRules, [][]int, error) {
	blocks, err := scanner.ReadBlocks(r)
	if err != nil {
		return nil, nil, err
	}
	if len(blocks) != 2 || len(blocks[1]) == 0 {
		return nil, nil, fmt.Errorf("want page rules and page updates separated by one blank line")
	}

	rules := make(pageOrderRules)
	for i, text := range blocks[0] {
		left, right, ok := strings.Cut(text, "|")
		if !ok || strings.Contains(right, "|") {
			return nil, nil, fmt.Errorf("line %d: invalid page rule %q", i+1, text)
		}
		before, err := pageOrderNumber(left)
		if err != nil {
			return nil, nil, fmt.Errorf("line %d: page rule: %w", i+1, err)
		}
		after, err := pageOrderNumber(right)
		if err != nil {
			return nil, nil, fmt.Errorf("line %d: page rule: %w", i+1, err)
		}
		rules[[2]int{before, after}] = true
	}

	updates := make([][]int, 0, len(blocks[1]))
	for i, text := range blocks[1] {
		line := len(blocks[0]) + 2 + i
		var update []int
		for field := range strings.SplitSeq(text, ",") {
			page, err := pageOrderNumber(field)
			if err != nil {
				return nil, nil, fmt.Errorf("line %d: page update: %w", line, err)
			}
			if slices.Contains(update, page) {
				return nil, nil, fmt.Errorf("line %d: page update repeats page %d", line, page)
			}
			update = append(update, page)
		}
		if len(update)%2 == 0 {
			return nil, nil, fmt.Errorf("line %d: page update needs an odd number of pages", line)
		}
		updates = append(updates, update)
	}
	return rules, updates, nil
}

func pageOrderNumber(field string) (int, error) {
	page, err := convert.ScanNumber[int]([]byte(field))
	if err != nil {
		return 0, err
	}
	if page < 0 {
		return 0, fmt.Errorf("page number %d is negative", page)
	}
	return page, nil
}

func pageOrderMiddleTotal(updates [][]int) (string, error) {
	var total int
	for _, update := range updates {
		sum, ok := xmath.CheckedAdd(total, update[len(update)/2])
		if !ok {
			return "", fmt.Errorf("middle page total overflows int")
		}
		total = sum
	}
	return strconv.Itoa(total), nil
}

func day05p01(r io.Reader) (string, error) {
	rules, updates, err := parsePageOrder(r)
	if err != nil {
		return "", err
	}
	var ordered [][]int
	for _, update := range updates {
		if slices.IsSortedFunc(update, rules.compare) {
			ordered = append(ordered, update)
		}
	}
	return pageOrderMiddleTotal(ordered)
}

func day05p02(r io.Reader) (string, error) {
	rules, updates, err := parsePageOrder(r)
	if err != nil {
		return "", err
	}
	var reordered [][]int
	for _, update := range updates {
		if !slices.IsSortedFunc(update, rules.compare) {
			slices.SortFunc(update, rules.compare)
			reordered = append(reordered, update)
		}
	}
	return pageOrderMiddleTotal(reordered)
}
