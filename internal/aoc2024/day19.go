package aoc2024

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/funcs"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
	"github.com/jacoelho/advent-of-code-go/pkg/xslices"
	"github.com/jacoelho/advent-of-code-go/pkg/xstrings"
)

func parseDesigns(r io.Reader) (collections.Set[string], []string, error) {
	const inputLimit = 128 << 10
	data, err := io.ReadAll(io.LimitReader(r, inputLimit+1))
	if err != nil {
		return nil, nil, err
	}
	if len(data) > inputLimit {
		return nil, nil, fmt.Errorf("towel input exceeds 128 KiB")
	}
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	patternsString, desiredString, ok := strings.Cut(text, "\n\n")
	if !ok {
		return nil, nil, fmt.Errorf("towels need patterns and desired designs")
	}

	patterns := collections.NewSet[string]()
	for field := range strings.SplitSeq(patternsString, ",") {
		pattern := strings.TrimSpace(field)
		if pattern == "" {
			return nil, nil, fmt.Errorf("towel pattern is empty")
		}
		patterns.Add(pattern)
	}
	return patterns, strings.Fields(desiredString), nil
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
	longestDesign := xiter.Max(xiter.Map(func(s string) int { return len(s) }, designs.Iter()))

	matches := xslices.CountFunc(func(s string) bool {
		return countMatches(designs, longestDesign, s) > 0
	}, patterns)
	return strconv.Itoa(matches), nil
}

func day19p02(r io.Reader) (string, error) {
	designs, patterns, err := parseDesigns(r)
	if err != nil {
		return "", err
	}
	longestDesign := xiter.Max(xiter.Map(func(s string) int { return len(s) }, designs.Iter()))

	count := xslices.Reduce(func(sum int, pattern string) int {
		return sum + countMatches(designs, longestDesign, pattern)
	}, 0, patterns)

	return strconv.Itoa(count), nil
}
