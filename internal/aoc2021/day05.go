package aoc2021

import (
	"bufio"
	"cmp"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

type ventGeometry byte

const (
	axisAlignedVents ventGeometry = iota
	diagonalVents
)

func ventPoint(text string) (grid.Position2D[int], error) {
	point, err := grid.ParsePosition2D[int](text)
	if err != nil {
		return point, err
	}
	if point.X < 0 || point.Y < 0 {
		return point, fmt.Errorf("negative vent point %q", text)
	}
	return point, nil
}

func ventOverlaps(r io.Reader, geometry ventGeometry) (int, error) {
	s := bufio.NewScanner(r)
	coverage := make(map[grid.Position2D[int]]int)
	overlaps := 0
	for s.Scan() {
		left, right, ok := strings.Cut(s.Text(), " -> ")
		if !ok {
			return 0, fmt.Errorf("invalid vent line %q", s.Text())
		}
		start, err := ventPoint(left)
		if err != nil {
			return 0, err
		}
		end, err := ventPoint(right)
		if err != nil {
			return 0, err
		}
		dx, dy := xmath.Abs(end.X-start.X), xmath.Abs(end.Y-start.Y)
		if dx != 0 && dy != 0 {
			if dx != dy {
				return 0, fmt.Errorf("vent line is not horizontal, vertical, or diagonal")
			}
			if geometry == axisAlignedVents {
				continue
			}
		}
		for step := 0; step <= max(dx, dy); step++ {
			point := grid.Position2D[int]{X: start.X + step*cmp.Compare(end.X, start.X), Y: start.Y + step*cmp.Compare(end.Y, start.Y)}
			coverage[point]++
			if coverage[point] == 2 {
				overlaps++
			}
		}
	}
	return overlaps, s.Err()
}

func day05p01(r io.Reader) (string, error) {
	overlaps, err := ventOverlaps(r, axisAlignedVents)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(overlaps), nil
}

func day05p02(r io.Reader) (string, error) {
	overlaps, err := ventOverlaps(r, diagonalVents)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(overlaps), nil
}
