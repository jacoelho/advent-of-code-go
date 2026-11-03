package aoc2024

import (
	"bytes"
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/search"
)

func parseFallingBytes(r io.Reader, dimensions, steps int) ([]grid.Position2D[int], error) {
	if dimensions < 0 || dimensions > 1024 {
		return nil, fmt.Errorf("memory dimensions must be between 0 and 1024")
	}
	count := 0
	s := scanner.NewScanner(r, func(line []byte) (grid.Position2D[int], error) {
		count++
		if count > 100_000 {
			return grid.Position2D[int]{}, fmt.Errorf("memory input exceeds 100000 falling bytes")
		}
		xText, yText, ok := bytes.Cut(line, []byte(","))
		if !ok {
			return grid.Position2D[int]{}, fmt.Errorf("falling byte needs x,y coordinates")
		}
		x, err := convert.ScanNumber[int](xText)
		if err != nil {
			return grid.Position2D[int]{}, fmt.Errorf("falling byte x: %w", err)
		}
		y, err := convert.ScanNumber[int](yText)
		if err != nil {
			return grid.Position2D[int]{}, fmt.Errorf("falling byte y: %w", err)
		}
		if x < 0 || x > dimensions || y < 0 || y > dimensions {
			return grid.Position2D[int]{}, fmt.Errorf("falling byte %d,%d is outside memory", x, y)
		}
		return grid.NewPosition2D(x, y), nil
	})
	corrupted := slices.Collect(s.Values())
	if err := s.Err(); err != nil {
		return nil, err
	}
	if steps < 0 || steps > len(corrupted) {
		return nil, fmt.Errorf("memory prefix %d is outside 0..%d", steps, len(corrupted))
	}
	return corrupted, nil
}

func memoryEdges(memory map[grid.Position2D[int]]bool, dimensions int) search.Edges[grid.Position2D[int]] {
	return func(p grid.Position2D[int], yield func(grid.Position2D[int], int) bool) error {
		for next := range grid.Neighbours4(p) {
			if next.X <= dimensions && next.Y <= dimensions && !memory[next] && !yield(next, 1) {
				return nil
			}
		}
		return nil
	}
}

func memoryHeuristic(dimensions int) func(grid.Position2D[int]) int {
	return func(p grid.Position2D[int]) int {
		return 2*dimensions - p.X - p.Y
	}
}

func memoryRoute(corrupted []grid.Position2D[int], dimensions, steps int) (int, bool, error) {
	memory := make(map[grid.Position2D[int]]bool, steps)
	for _, position := range corrupted[:steps] {
		memory[position] = true
	}
	start := grid.NewPosition2D(0, 0)
	end := grid.NewPosition2D(dimensions, dimensions)
	if memory[start] || memory[end] {
		return 0, false, nil
	}
	return search.AStarDistance(
		start,
		func(position grid.Position2D[int]) bool { return position == end },
		memoryEdges(memory, dimensions),
		memoryHeuristic(dimensions),
		search.Limits{},
	)
}

func day18p01(dimensions, steps int) func(io.Reader) (string, error) {
	return func(r io.Reader) (string, error) {
		corrupted, err := parseFallingBytes(r, dimensions, steps)
		if err != nil {
			return "", err
		}
		cost, found, err := memoryRoute(corrupted, dimensions, steps)
		if err != nil {
			return "", fmt.Errorf("memory route: %w", err)
		}
		if !found {
			return "", fmt.Errorf("no route through memory after %d falling bytes", steps)
		}
		return strconv.Itoa(cost), nil
	}
}

func day18p02(dimensions, steps int) func(io.Reader) (string, error) {
	return func(r io.Reader) (string, error) {
		corrupted, err := parseFallingBytes(r, dimensions, steps)
		if err != nil {
			return "", err
		}
		if _, found, err := memoryRoute(corrupted, dimensions, steps); err != nil {
			return "", fmt.Errorf("memory prefix route: %w", err)
		} else if !found {
			return "", fmt.Errorf("memory prefix of %d bytes already blocks the route", steps)
		}
		if steps == len(corrupted) {
			return "", fmt.Errorf("no falling byte blocks the route")
		}
		prefix, err := search.BinarySearch(steps+1, len(corrupted), func(prefix int) (bool, error) {
			_, found, err := memoryRoute(corrupted, dimensions, prefix)
			return !found, err
		})
		if err != nil {
			return "", fmt.Errorf("find first blocking byte: %w", err)
		}
		position := corrupted[prefix-1]
		return fmt.Sprintf("%d,%d", position.X, position.Y), nil
	}
}
