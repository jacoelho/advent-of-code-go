package aoc2023

import (
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
)

type springRecord struct {
	pattern string
	groups  []int
}

func parseSpringRecord(line []byte) (springRecord, error) {
	parts := strings.Fields(string(line))
	if len(parts) != 2 {
		return springRecord{}, nil
	}

	return springRecord{
		pattern: parts[0],
		groups:  convert.ExtractDigits[int](parts[1]),
	}, nil
}

func countArrangements(pattern string, groups []int) int {
	cache := make(map[[3]int]int)

	var count func(pos, groupIdx, currentRun int) int
	count = func(pos, groupIdx, currentRun int) int {
		key := [3]int{pos, groupIdx, currentRun}
		if cached, found := cache[key]; found {
			return cached
		}

		if pos == len(pattern) {
			if groupIdx == len(groups) && currentRun == 0 {
				return 1
			}
			if groupIdx == len(groups)-1 && currentRun == groups[groupIdx] {
				return 1
			}
			return 0
		}

		result := 0
		ch := pattern[pos]

		if ch == '.' || ch == '?' {
			if currentRun == 0 {
				result += count(pos+1, groupIdx, 0)
			} else if groupIdx < len(groups) && currentRun == groups[groupIdx] {
				result += count(pos+1, groupIdx+1, 0)
			}
		}

		if ch == '#' || ch == '?' {
			result += count(pos+1, groupIdx, currentRun+1)
		}

		cache[key] = result
		return result
	}

	return count(0, 0, 0)
}

func unfold(rec springRecord, times int) springRecord {
	patterns := make([]string, times)
	for i := range patterns {
		patterns[i] = rec.pattern
	}

	groups := make([]int, 0, len(rec.groups)*times)
	for range times {
		groups = append(groups, rec.groups...)
	}

	return springRecord{
		pattern: strings.Join(patterns, "?"),
		groups:  groups,
	}
}

func solve(unfoldFactor int) func(io.Reader) (string, error) {
	return func(r io.Reader) (string, error) {
		s := scanner.NewScanner(r, parseSpringRecord)
		total := xiter.Sum(xiter.Map(func(rec springRecord) int {
			if unfoldFactor > 1 {
				rec = unfold(rec, unfoldFactor)
			}
			return countArrangements(rec.pattern, rec.groups)
		}, s.Values()))
		if err := s.Err(); err != nil {
			return "", err
		}
		return strconv.Itoa(total), nil
	}
}

func day12p01(r io.Reader) (string, error) {
	return solve(1)(r)
}

func day12p02(r io.Reader) (string, error) {
	return solve(5)(r)
}
