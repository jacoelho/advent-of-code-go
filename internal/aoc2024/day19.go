package aoc2024

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/funcs"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
	"github.com/jacoelho/advent-of-code-go/pkg/xstrings"
)

func parseDesigns(r io.Reader) (collections.Set[string], []string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, nil, err
	}
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	patternsString, desiredString, ok := strings.Cut(text, "\n\n")
	if !ok {
		return nil, nil, fmt.Errorf("towels need patterns and desired designs")
	}

	patterns := collections.NewSet[string]()
	for field := range strings.SplitSeq(patternsString, ",") {
		pattern := strings.TrimSpace(field)
		if err := day19ValidStripes(pattern); err != nil {
			return nil, nil, fmt.Errorf("towel pattern: %w", err)
		}
		patterns.Add(pattern)
	}
	designs := strings.Fields(desiredString)
	for _, design := range designs {
		if err := day19ValidStripes(design); err != nil {
			return nil, nil, fmt.Errorf("towel design: %w", err)
		}
	}
	return patterns, designs, nil
}

func day19ValidStripes(stripes string) error {
	if stripes == "" {
		return fmt.Errorf("no stripes")
	}
	if i := strings.IndexFunc(stripes, func(c rune) bool { return !strings.ContainsRune("wubrg", c) }); i >= 0 {
		return fmt.Errorf("%q has invalid stripe colour %q", stripes, stripes[i])
	}
	return nil
}

func countMatches(designs collections.Set[string], longestDesign int, pattern string) int {
	var countMatchesMemoized func(pattern string) int
	countMatchesMemoized = funcs.Memoize(func(s string) int {
		if s == "" {
			return 1
		}
		var count int
		for _, p := range xstrings.SubSlices(s, longestDesign) {
			if designs.Contains(p[0]) {
				count += countMatchesMemoized(p[1])
			}
		}
		return count
	})
	return countMatchesMemoized(pattern)
}

func day19p01(r io.Reader) (string, error) {
	designs, patterns, err := parseDesigns(r)
	if err != nil {
		return "", err
	}
	longestDesign, _ := xiter.Max(xiter.From(designs.Iter()).Map(func(s string) int { return len(s) }).All())

	matches := xiter.FromSlice(patterns).Filter(func(s string) bool {
		return countMatches(designs, longestDesign, s) > 0
	}).Count()
	return strconv.Itoa(matches), nil
}

func day19p02(r io.Reader) (string, error) {
	designs, patterns, err := parseDesigns(r)
	if err != nil {
		return "", err
	}
	longestDesign, _ := xiter.Max(xiter.From(designs.Iter()).Map(func(s string) int { return len(s) }).All())

	count := xiter.FromSlice(patterns).SumBy(func(pattern string) int {
		return countMatches(designs, longestDesign, pattern)
	})

	return strconv.Itoa(count), nil
}
