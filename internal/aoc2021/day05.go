package aoc2021

import (
	"bufio"
	"cmp"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

const maxVentPoints = 1_000_000

type ventGeometry byte

const (
	axisAlignedVents ventGeometry = iota
	diagonalVents
)

func ventPoint(text string) (grid.Position2D[int], error) {
	var point grid.Position2D[int]
	x, y, ok := strings.Cut(text, ",")
	if !ok {
		return point, fmt.Errorf("invalid vent point %q", text)
	}
	var err error
	point.X, err = convert.ScanNumber[int]([]byte(strings.TrimSpace(x)))
	if err != nil {
		return point, err
	}
	point.Y, err = convert.ScanNumber[int]([]byte(strings.TrimSpace(y)))
	return point, err
}

func ventDistance(a, b int) uint {
	if a < b {
		return uint(b) - uint(a)
	}
	return uint(a) - uint(b)
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
		dx, dy := ventDistance(start.X, end.X), ventDistance(start.Y, end.Y)
		if dx != 0 && dy != 0 {
			if dx != dy {
				return 0, fmt.Errorf("vent line is not horizontal, vertical, or diagonal")
			}
			if geometry == axisAlignedVents {
				continue
			}
		}
		steps := max(dx, dy)
		if steps >= maxVentPoints {
			return 0, fmt.Errorf("vent line exceeds %d raster points", maxVentPoints)
		}
		for step := range int(steps) + 1 {
			point := grid.Position2D[int]{X: start.X + step*cmp.Compare(end.X, start.X), Y: start.Y + step*cmp.Compare(end.Y, start.Y)}
			if coverage[point] == 0 && len(coverage) == maxVentPoints {
				return 0, fmt.Errorf("vent coverage exceeds %d points", maxVentPoints)
			}
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
