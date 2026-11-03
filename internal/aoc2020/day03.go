package aoc2020

import (
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

func parseTreeMap(r io.Reader) ([]string, error) {
	width := 0
	s := scanner.NewScanner(r, func(line []byte) (string, error) {
		if len(line) == 0 {
			return "", fmt.Errorf("empty tree map row")
		}
		if width == 0 {
			width = len(line)
		}
		if len(line) != width {
			return "", fmt.Errorf("tree map row has width %d, want %d", len(line), width)
		}
		for _, c := range line {
			if c != '.' && c != '#' {
				return "", fmt.Errorf("invalid tree map character %q", c)
			}
		}
		return string(line), nil
	})
	rows := slices.Collect(s.Values())
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("empty tree map")
	}
	return rows, nil
}

func treesOnSlope(rows []string, right, down int) int64 {
	var trees int64
	for x, y := right, down; y < len(rows); x, y = x+right, y+down {
		if rows[y][x%len(rows[y])] == '#' {
			trees++
		}
	}
	return trees
}

func day03p01(r io.Reader) (string, error) {
	rows, err := parseTreeMap(r)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(treesOnSlope(rows, 3, 1), 10), nil
}

func day03p02(r io.Reader) (string, error) {
	rows, err := parseTreeMap(r)
	if err != nil {
		return "", err
	}
	product := int64(1)
	for _, slope := range [][2]int{{1, 1}, {3, 1}, {5, 1}, {7, 1}, {1, 2}} {
		product *= treesOnSlope(rows, slope[0], slope[1])
	}
	return strconv.FormatInt(product, 10), nil
}
