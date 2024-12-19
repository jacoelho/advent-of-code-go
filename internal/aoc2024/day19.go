package aoc2024

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/internal/collections"
	"github.com/jacoelho/advent-of-code-go/internal/funcs"
	"github.com/jacoelho/advent-of-code-go/internal/scanner"
	"github.com/jacoelho/advent-of-code-go/internal/xiter"
	"github.com/jacoelho/advent-of-code-go/internal/xslices"
	"github.com/jacoelho/advent-of-code-go/internal/xstrings"
)

func parseDesigns(r io.Reader) (collections.Set[string], []string) {
	s := bufio.NewScanner(r)
	s.Split(scanner.SplitBySeparator([]byte{'\n', '\n'}))
	_, patternsString := s.Scan(), s.Text()
	_, desiredString := s.Scan(), s.Text()
	if s.Err() != nil {
		panic(s.Err())
	}

	patterns := collections.NewSet(strings.Split(patternsString, ", ")...)
	desired := strings.Fields(desiredString)

	return patterns, desired
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
	designs, patterns := parseDesigns(r)
	longestDesign := xiter.Max(xiter.Map(func(s string) int { return len(s) }, designs.Iter()))

	matches := xslices.CountFunc(func(s string) bool {
		return countMatches(designs, longestDesign, s) > 0
	}, patterns)
	return strconv.Itoa(matches), nil
}

func day19p02(r io.Reader) (string, error) {
	designs, patterns := parseDesigns(r)
	longestDesign := xiter.Max(xiter.Map(func(s string) int { return len(s) }, designs.Iter()))

	count := xslices.Reduce(func(sum int, pattern string) int {
		return sum + countMatches(designs, longestDesign, pattern)
	}, 0, patterns)

	return strconv.Itoa(count), nil
}
