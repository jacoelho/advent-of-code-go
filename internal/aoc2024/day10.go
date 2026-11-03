package aoc2024

import (
	"io"
	"iter"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/graph"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/search"
)

func parseTopographicMap(r io.Reader) ([][]byte, error) {
	s := scanner.NewScanner(r, func(b []byte) ([]byte, error) {
		res := make([]byte, len(b))
		// Preserve the byte-wrapped heights accepted for nondigits.
		for i, v := range b {
			res[i] = v - '0'
		}
		return res, nil
	})
	return slices.Collect(s.Values()), s.Err()
}

func day10neighbours(heights [][]byte) func(grid.Position2D[int], func(grid.Position2D[int]) bool) {
	return func(p grid.Position2D[int], yield func(grid.Position2D[int]) bool) {
		slope := int(heights[p.Y][p.X])
		for _, offset := range grid.OffsetsNeighbours4[int]() {
			v := p.Add(offset)
			if v.Y < 0 || v.Y >= len(heights) || v.X < 0 || v.X >= len(heights[v.Y]) {
				continue
			}
			if int(heights[v.Y][v.X]) == slope+1 && !yield(v) {
				return
			}
		}
	}
}

func trailheads(heights [][]byte) []grid.Position2D[int] {
	var trailheads []grid.Position2D[int]
	for y, row := range heights {
		for x, slope := range row {
			if slope == 0 {
				trailheads = append(trailheads, grid.Position2D[int]{X: x, Y: y})
			}
		}
	}
	return trailheads
}

func day10p01(r io.Reader) (string, error) {
	heights, err := parseTopographicMap(r)
	if err != nil {
		return "", err
	}

	visit := day10neighbours(heights)
	neighbours := func(p grid.Position2D[int]) iter.Seq[grid.Position2D[int]] {
		return func(yield func(grid.Position2D[int]) bool) { visit(p, yield) }
	}
	var summits int
	for _, trail := range trailheads(heights) {
		for position := range search.BFS(trail, neighbours) {
			if heights[position.Y][position.X] == 9 {
				summits++
			}
		}
	}

	return strconv.Itoa(summits), nil
}

func day10p02(r io.Reader) (string, error) {
	heights, err := parseTopographicMap(r)
	if err != nil {
		return "", err
	}

	summits, err := graph.CountDAGPaths(trailheads(heights), day10neighbours(heights), func(position grid.Position2D[int]) (int, bool) {
		return 1, heights[position.Y][position.X] == 9
	})
	if err != nil {
		return "", err
	}
	return strconv.Itoa(summits), nil
}
