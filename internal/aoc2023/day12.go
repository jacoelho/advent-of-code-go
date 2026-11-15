package aoc2023

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
)

type springRecord struct {
	pattern string
	groups  []int
}

func parseSpringRecord(line []byte) (springRecord, error) {
	pattern, groupsText, ok := strings.Cut(string(line), " ")
	if !ok || pattern == "" || strings.Trim(pattern, ".#?") != "" {
		return springRecord{}, fmt.Errorf("invalid spring record %q", line)
	}
	var groups []int
	for field := range strings.SplitSeq(groupsText, ",") {
		group, err := strconv.Atoi(field)
		if err != nil || group < 1 {
			return springRecord{}, fmt.Errorf("invalid damaged group %q in %q", field, line)
		}
		groups = append(groups, group)
	}
	return springRecord{pattern: pattern, groups: groups}, nil
}

func countArrangements(pattern string, groups []int) int {
	n, g := len(pattern), len(groups)
	// operational[i] is the length of the run of non-operational springs at i.
	operational := make([]int, n+1)
	for i := n - 1; i >= 0; i-- {
		if pattern[i] != '.' {
			operational[i] = operational[i+1] + 1
		}
	}
	// ways[i*(g+1)+j] counts arrangements of groups[j:] in pattern[i:]. Row n+1
	// mirrors row n so a group ending at n can skip its separator uniformly.
	stride := g + 1
	ways := make([]int, (n+2)*stride)
	ways[n*stride+g], ways[(n+1)*stride+g] = 1, 1
	for i := n - 1; i >= 0; i-- {
		for j := 0; j <= g; j++ {
			count := 0
			if pattern[i] != '#' {
				count += ways[(i+1)*stride+j]
			}
			if pattern[i] != '.' && j < g {
				end := i + groups[j]
				if operational[i] >= groups[j] && (end == n || pattern[end] != '#') {
					count += ways[(end+1)*stride+j+1]
				}
			}
			ways[i*stride+j] = count
		}
	}
	return ways[0]
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
	return springRecord{pattern: strings.Join(patterns, "?"), groups: groups}
}

func springArrangements(r io.Reader, unfoldFactor int) (string, error) {
	s := scanner.NewScanner(r, parseSpringRecord)
	total := xiter.From(s.Values()).SumBy(func(rec springRecord) int {
		rec = unfold(rec, unfoldFactor)
		return countArrangements(rec.pattern, rec.groups)
	})
	if err := s.Err(); err != nil {
		return "", err
	}
	return strconv.Itoa(total), nil
}

func day12p01(r io.Reader) (string, error) {
	return springArrangements(r, 1)
}

func day12p02(r io.Reader) (string, error) {
	return springArrangements(r, 5)
}
